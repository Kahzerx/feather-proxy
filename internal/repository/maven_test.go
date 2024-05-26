package repository

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMavenConnection(t *testing.T) {
	config := NewMavenConfig()
	get, err := http.Get(fmt.Sprintf("%s://%s:%s", config.Scheme, config.Host, config.Port))
	if err != nil {
		t.Fatal("Unable to connect to maven repository, err: ", err.Error())
	}
	if get.StatusCode != http.StatusOK {
		t.Fatal("Unable to connect to maven repository, status code: ", get.StatusCode)
	}
}
