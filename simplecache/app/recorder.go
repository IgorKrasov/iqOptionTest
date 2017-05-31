package app

import (
	"time"
	"os"
	"io"
	"encoding/gob"
	"fmt"
	"sync"
)

type recorder struct {
	Interval time.Duration
	mu                    sync.RWMutex
	fname    string
	stop chan bool
}

func (r *recorder) run (c *cache) {
	r.stop = make(chan bool)
	ticker := time.NewTicker(r.Interval)
	for {
		select {
		case <- ticker.C:
			r.saveToFile(c.Items())
		case <- r.stop:
			ticker.Stop()
			return
		}
	}
}

func(r *recorder) saveToFile(items map[string]item) error {
	println("save to file")
	fp, err := os.Create(r.fname)
	if err != nil {
		return err
	}
	err = r.Save(fp, items)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (r *recorder) Save(w io.Writer, items map[string]item) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, vi := range items {
		switch vi.(type) {
		case simpleItem:
			//v := vi.(simpleItem)
			//gob.Register(v.object)
		case listItem:
			//v := vi.(listItem)
			gob.Register([]interface{}{})
		case dictItem:
			//v := vi.(dictItem)
			gob.Register(map[string]interface{}{})
		}
	}
	err = enc.Encode(&items)
	if err != nil {
		println(err.Error())
	}
	return
}

//func stopRecorder(c *cache) {
//	c.recorder.stop <- true
//}

func runRecorder(c* cache, ri time.Duration) {
	r := &recorder{
		Interval: ri,
		fname: "test.txt",
	}
	c.recorder = r
	go r.run(c)
}
