package main

import (
	"JenkinsUpdateAgent/src/JenkinsUpdateAgent"
	"fmt"
	"log"
	"net/http"
)

const (
	servicePort = 8888
	serviceIP   = "0.0.0.0"
)

func main() {
	log.Printf("Service Start %s:%d ...\n", serviceIP, servicePort)
	http.HandleFunc("/", JenkinsUpdateAgent.HandleDefault)
	http.HandleFunc("/updates/", JenkinsUpdateAgent.HandleOtherUpdate)
	http.HandleFunc("/update-center.json", JenkinsUpdateAgent.HandleUpdateJson)
	_ = http.ListenAndServe(fmt.Sprintf("%s:%d", serviceIP, servicePort), nil)
}
