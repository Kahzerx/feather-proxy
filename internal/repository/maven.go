package repository

import "os"

type MavenConfig struct {
	Scheme string
	Host   string
	Port   string
}

func NewMavenConfig() *MavenConfig {
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
	return &MavenConfig{
		Scheme: mavenScheme,
		Host:   mavenHost,
		Port:   mavenPort,
	}
}
