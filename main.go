package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const htmlPrefix = "<html><body><h2>Here</h2><hr><ul>\n"
const htmlPostfix = "</ul></body></html>"

func requestHandler(w http.ResponseWriter, req *http.Request) {

	filePath := filepath.Join(".", req.RequestURI)
	fmt.Println(filePath)

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

func run(portNumber int) error {

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	portNumberStr := strconv.Itoa(portNumber)
	fmt.Println("[Here] - serving " + path + " at port " + portNumberStr)
	http.HandleFunc("/", requestHandler)
	if err := http.ListenAndServe("127.0.0.1:"+portNumberStr, nil); err != nil {
		return err
	}
	return nil
}
func main() {

	var portNumber int

	flag.IntVar(&portNumber, "port", 9898, "port number to listen on")
	flag.Parse()

	if err := run(portNumber); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
