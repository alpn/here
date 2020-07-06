package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const htmlPrefix = "<html><body><h2>Here</h2><hr><ul>\n"
const htmlPostfix = "</ul></body></html>"

func requestHandler(w http.ResponseWriter, req *http.Request) {

	filePath := "./" + req.RequestURI

	fileStat, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - File Not Found"))
		return
	}

	if nil != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	if fileStat.IsDir() {
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Fatal(err)
		}

		var html = htmlPrefix

		for _, file := range files {
			var line = file.Name()
			if file.IsDir() {
				line += "/"
			}
			html += "<li><a href='" + line + "'>" + line + "</a>\n"
		}

		html += htmlPostfix
		buffer := []byte(html)
		w.Write(buffer)
		return

	}

	f, err := os.Open(filePath)
	defer f.Close()
	if _, err = io.Copy(w, f); err != nil {
		fmt.Println(err)
	}
}

func main() {

	http.HandleFunc("/", requestHandler)
	if err := http.ListenAndServe(":9898", nil); err != nil {
		fmt.Println(err)
	}
}
