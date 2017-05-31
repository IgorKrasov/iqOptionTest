package main

import (
	"github.com/iqOptionTest/simplecacheclient/app"
	"flag"
	"net/http"
	"fmt"
)

const ServerUrl = "http://cache:9003/"

var (
        loadbalancerURL = flag.String("loadbalancer", "http://balancer:9005", "Address of the load balancer")
	appName = "cacheClient"
)

func main() {
	a := app.New(appName, ":9004", ServerUrl)
	fmt.Println("run server")
	http.Get(*loadbalancerURL + "/register?port=9004")
	a.Run()
	fmt.Scanln()
	http.Get(*loadbalancerURL + "/register?port=9004")
}
