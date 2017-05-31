package main

import (
	"net/http"
	"fmt"
	"strings"
	"time"
	"io"
	"net/url"
)

type webRequest struct {
	r *http.Request
	w http.ResponseWriter
	doneCh chan struct{}
}

var (
	requestCh = make(chan *webRequest)
	registerCh = make(chan string)
	unregisterCh = make(chan string)
	healthCheckCh = time.Tick(5 * time.Second)
	appservers = []string{}
	currentIndex = 0
	client = http.Client{}
)

func appServerRegister(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	println(r.RemoteAddr)
	port := r.URL.Query().Get("port")
	println(ip)
	registerCh <- ip + ":" + port
}

func appServerUnregister(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	println(r.RemoteAddr)
	port := r.URL.Query().Get("port")
	println(ip)
	unregisterCh <- ip + ":" + port
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})
		requestCh <- &webRequest{
			r: r,
			w: w,
			doneCh: doneCh,
		}
		<-doneCh
	})
	http.HandleFunc("/register", appServerRegister)
	http.HandleFunc("/unregister", appServerUnregister)

	go ProcessRequests()
	http.ListenAndServe(":9005", nil)
	fmt.Println("start server on port :9005")
	fmt.Scanln()
}

func ProcessRequests() {
	for {
		select {
		case request := <-requestCh:
			println("request")
			if len(appservers) == 0 {
				request.w.WriteHeader(http.StatusInternalServerError)
				request.w.Write([]byte("No app servers found"))
				request.doneCh <- struct{}{}
				continue
			}
			//currentIndex++
			if currentIndex > len(appservers) - 1 {
				currentIndex = 0
			}
			host := appservers[currentIndex]

			go processRequest(host, request)
		case host := <-registerCh:
			println("register ", host)
			isFound := false
			for _, h := range appservers {
				if host == h {
					isFound = true
					break
				}
			}
			if !isFound {
				appservers = append(appservers, host)
			}
		case host := <-unregisterCh:
			println("unregister", host)
			for i := len(appservers) - 1; i >= 0; i-- {
				if appservers[i] == host {
					appservers = append(appservers[:i], appservers[i+1:]...)
				}
			}
		case <-healthCheckCh:
			println("healthCheck ")
			servers := appservers[:]
			go func(servers []string) {
				for _, server := range servers {
					resp, err := http.Get("http://" + server + "/ping")
					if err != nil || resp.StatusCode != http.StatusOK {
						unregisterCh <- server
					}
				}
			}(servers)
		}
	}
}

func processRequest(host string, request *webRequest) {
	hostURL, _ := url.Parse(request.r.URL.String())
	hostURL.Scheme = "http"
	hostURL.Host = host
	println(host)
	println(hostURL.String())
	req, _ := http.NewRequest(request.r.Method, hostURL.String(), request.r.Body)
	for k, v := range request.r.Header {
		if k == "Content-Length" || k == "Last-Modified" {
			continue
		}
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		req.Header.Add(k, values)
	}

	resp, err := client.Do(req)

	if err != nil {
		request.w.WriteHeader(http.StatusInternalServerError)
		request.doneCh <- struct{}{}
		return
	}

	for k, v := range resp.Header {
		println(k)
		if k == "Content-Length" || k == "Last-Modified" {
			continue
		}
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		request.w.Header().Add(k, values)
	}
	io.Copy(request.w, resp.Body)

	request.doneCh <- struct{}{}
}
