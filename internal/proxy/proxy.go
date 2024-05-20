package proxy

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type proxy struct {
	mutex *redsync.Mutex
}

func NewProxy() http.Handler {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rs := redsync.New(goredis.NewPool(client))
	mutex := rs.NewMutex("feather-publishing", redsync.WithExpiry(10*time.Second))
	return &proxy{
		mutex: mutex,
	}
}

func (p *proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println(request.Method + " " + request.URL.String())
	if request.Method == http.MethodGet && strings.HasSuffix(request.URL.String(), "maven-metadata.xml") {
		_ = p.mutex.Lock()
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
		_, _ = p.mutex.Unlock()
	}
}
