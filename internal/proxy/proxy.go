package proxy

import (
	"context"
	"feather-proxy/internal/database"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type mavenConfig struct {
	scheme string
	host   string
	port   string
}

type proxy struct {
	mutex  *redsync.Mutex
	config mavenConfig
}

func NewProxy() http.Handler {
	client := database.NewRedisClient()
	if client.Ping(context.Background()).Err() != nil {
		panic("Unable to connect to redis")
	}
	rs := redsync.New(goredis.NewPool(client))
	mutex := rs.NewMutex("feather-publishing", redsync.WithExpiry(10*time.Second))
	mavenScheme := os.Getenv("MAVEN_SCHEME")
	if mavenScheme == "" {
		mavenScheme = "http"
	}
	mavenHost := os.Getenv("MAVEN_HOST")
	if mavenHost == "" {
		mavenHost = "127.0.0.1"
	}
	mavenPort := os.Getenv("MAVEN_PORT")
	if mavenPort == "" {
		mavenPort = "80"
	}
	return &proxy{
		mutex: mutex,
		config: mavenConfig{
			scheme: mavenScheme,
			host:   mavenHost,
			port:   mavenPort,
		},
	}
}

func (p *proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println(request.Method + " " + request.URL.String())
	if request.Method == http.MethodGet && strings.HasSuffix(request.URL.String(), "maven-metadata.xml") {
		_ = p.mutex.Lock()
	}
	client := &http.Client{}
	request.URL.Host = fmt.Sprintf("%s:%s", p.config.host, p.config.port)
	request.Host = fmt.Sprintf("%s:%s", p.config.host, p.config.port)
	request.URL.Scheme = p.config.scheme
	request.RequestURI = ""
	// Comment if reverse proxy already does the redirect
	// request.URL.Path = "/releases" + request.URL.Path
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
	if request.Method == http.MethodPut && strings.HasSuffix(request.URL.String(), "maven-metadata.xml.sha512") {
		_, _ = p.mutex.Unlock()
	}
}
