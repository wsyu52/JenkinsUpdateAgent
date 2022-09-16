package JenkinsUpdateAgent

import (
	"log"
	"net/http"
	"strings"
)

func showRequestInfo(r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
}

func updateJsonContent(body string, mirrorUrl string) string {
	jsonBody := strings.Replace(body, "https://www.google.com", "https://www.baidu.com", -2)
	jsonBody = strings.Replace(jsonBody, "https://updates.jenkins.io/download/plugins/", mirrorUrl+"/jenkins/plugins/", -2)
	return jsonBody
}
