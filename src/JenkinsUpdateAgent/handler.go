package JenkinsUpdateAgent

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	MirrorURL = "https://mirrors.tuna.tsinghua.edu.cn"
)

func HandleDefault(w http.ResponseWriter, r *http.Request) {
	showRequestInfo(r)
	_, _ = fmt.Fprintln(w, "hello from jenkins updates agent.")
}

func HandleOtherUpdate(w http.ResponseWriter, r *http.Request) {
	showRequestInfo(r)
	redirectUrl := fmt.Sprintf("%s/jenkins/updates%s", MirrorURL, r.URL)
	log.Printf("redirect to: %s\n", redirectUrl)
	w.WriteHeader(http.StatusFound)
	w.Header().Set("location", redirectUrl)
}

func HandleUpdateJson(w http.ResponseWriter, r *http.Request) {
	showRequestInfo(r)

	// handle test
	if r.Form != nil && r.Form.Has("uctest") {
		log.Printf("jenkins update center test")
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{}")
		return
	}

	jenkinsVersion := r.FormValue("version")

	// no version info
	if jenkinsVersion == "" {
		log.Printf("no jenkins version info")
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{}")
		return
	}

	// try to get json from original version
	originalVersionURL := fmt.Sprintf("%s/jenkins/updates/dynamic-stable-%s/update-center.json", MirrorURL, jenkinsVersion)
	log.Printf("try to get json from: %s\n", originalVersionURL)
	originalResp, err := http.Get(originalVersionURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer originalResp.Body.Close()
	originalBody, err := io.ReadAll(originalResp.Body)
	log.Printf("status: %d content-type: %s body-type: %T body-len: %d\n",
		originalResp.StatusCode, originalResp.Header.Get("Content-Type"), originalBody, len(originalBody))

	if originalResp.StatusCode == http.StatusOK && originalResp.Header.Get("Content-Type") == "application/json" {
		log.Println("original version json exists")
		// return json
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, updateJsonContent(string(originalBody), MirrorURL))
		return
	} else if originalResp.StatusCode == http.StatusNotFound {
		// get version list
		jsonListUrl := MirrorURL + "/jenkins/updates/"
		jsonListResp, err := soup.Get(jsonListUrl)
		if err != nil {
			log.Println(err)
			return
		}
		var curVer string
		for _, v := range soup.HTMLParse(jsonListResp).FindAll("a") {
			if strings.Contains(v.Attrs()["href"], "dynamic-stable-") {
				ver := strings.Replace(strings.Replace(v.Attrs()["href"], "dynamic-stable-", "", -20), "/", "", 1)
				if jenkinsVersion > ver {
					curVer = ver
				}
			}
		}
		//handle ver
		curJsonURL := fmt.Sprintf("%s/jenkins/updates/dynamic-stable-%s/update-center.json", MirrorURL, curVer)
		log.Println("match version: " + curVer + ", url: " + curJsonURL)
		matchVersionResp, err := http.Get(curJsonURL)
		if err != nil {
			log.Println(err)
			return
		}
		defer matchVersionResp.Body.Close()
		matchVersionBody, err := io.ReadAll(matchVersionResp.Body)
		log.Printf("status: %d content-type: %s body-type: %T body-len: %d\n",
			matchVersionResp.StatusCode, matchVersionResp.Header.Get("Content-Type"), matchVersionBody, len(matchVersionBody))

		if matchVersionResp.StatusCode == http.StatusOK && matchVersionResp.Header.Get("Content-Type") == "application/json" {
			log.Println("json file exists")
			// return json
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, updateJsonContent(string(matchVersionBody), MirrorURL))
			return
		}
	}
	// default version
	defaultJsonUrl := MirrorURL + "/jenkins/updates/update-center.json"
	log.Printf("cannot match jenkins version, redirect to: " + defaultJsonUrl)
	w.WriteHeader(http.StatusFound)
	w.Header().Set("location", defaultJsonUrl)

}
