package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/muaazsaleem/cyao"
	"github.com/sirupsen/logrus"
)

const CONFIG_FILE = "gopher.json"

func check(err error, msg string) {
	if err != nil {
		logrus.Fatalf(msg, err)
	}
}

func main() {
	f, err := os.Open(CONFIG_FILE)
	check(err, "couldn't open config file: "+CONFIG_FILE)

	story := map[string]cyoa.Arc{}

	jsonData, err := ioutil.ReadAll(f)

	err = json.Unmarshal(jsonData, &story)
	check(err, "failed to unmarshal json")

	serverMux := defaultMux()
	handler := serveStory(serverMux, story)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		logrus.Fatalln(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func serveStory(mux *http.ServeMux, story map[string]cyoa.Arc) http.Handler {

	for title := range story {
		fmt.Println(title)
		mux.Handle("/"+title, archHandler())
	}

	return mux
}

func archHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL
		_, err := fmt.Fprintf(w, "Serving the arc %s", path)
		if err != nil {
			logrus.Errorf("failed to write to response for request %v", r, err)
		}
	})
}
