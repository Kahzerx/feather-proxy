package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type proxy struct {
	lock      *sync.Mutex
	locked    *atomic.Bool
	lockTimer *time.Timer
}

func NewProxy() http.Handler {
	return &proxy{
		lock:      &sync.Mutex{},
		lockTimer: time.NewTimer(0),
		locked:    &atomic.Bool{},
	}
}

func (p *proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println(request.Method + " " + request.URL.String())
	if request.Method == http.MethodGet && strings.HasSuffix(request.URL.String(), "maven-metadata.xml") {
		p.lock.Lock()
		p.locked.Store(true)
		p.lockTimer.Reset(20 * time.Second)
		go func() {
			<-p.lockTimer.C
			if p.locked.Load() {
				p.locked.Store(false)
				p.lock.Unlock()
			}
		}()
	}
	client := &http.Client{}
	// TODO host as env
	request.URL.Host = "0.0.0.0:80"
	request.Host = "0.0.0.0:80"
	// TODO schema as env
	request.URL.Scheme = "http"
	request.RequestURI = ""
	// Comment if reverse proxy already does the redirect
	request.URL.Path = "/releases" + request.URL.Path
	resp, err := client.Do(request)
	if err != nil {
		http.Error(writer, "Server Error", http.StatusInternalServerError)
		log.Println("ServeHTTP:", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	writer.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(writer, resp.Body)
	// TODO time based unlock!
	if request.Method == http.MethodPut && strings.HasSuffix(request.URL.String(), "maven-metadata.xml.sha512") {
		if !p.lockTimer.Stop() {
			<-p.lockTimer.C
		}
		if p.locked.Load() {
			p.locked.Store(false)
			p.lock.Unlock()
		}
	}
}
