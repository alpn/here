package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/djherbis/times"
	"github.com/gomarkdown/markdown"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func htmlPrefix(title string, customHead string) string {
	return `<!DOCTYPE html>` + `<html><head><title>` + title + `</title>` + customHead + `</head><body>`
}

const htmlPostfix = `</ul></body></html>`
const css = `
        <style>
            html {
              max-width: 70ch;
              padding: 3em 1em;
              margin: auto;
              line-height: 1.75;
              font-size: 1.25em;
              color : white;
              background-color : black;
              font-family: 'Courier New', monospace;
            }

            img {
              max-width: 70ch;
            }
            a{
                text-decoration: none;
                color : red
            }
        </style>`

func requestHandler(w http.ResponseWriter, req *http.Request) {

	filePath := filepath.Join(".", req.RequestURI)
	log.Println(filePath)

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

		title := "Here"
		var html = htmlPrefix(title, "")
		html += `<h2>` + title + `</h2><hr><ul>`

		sort.SliceStable(files, func(i, j int) bool {
			if files[i].IsDir() && !files[j].IsDir() {
				return true
			}
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		for _, file := range files {

			var line = file.Name()

			if line[0] == '.' {
				continue
			}

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

	// Render Markdown files to HTML
	if path.Ext(filePath) == ".md" {
		handleMarkdownFile(filePath, w)
		return
	}

	f, err := os.Open(filePath)
	defer f.Close()

	if path.Ext(filePath) == ".css" {
		w.Header().Add("Content-Type", "text/css ; charset=utf-8")
	}

	if _, err = io.Copy(w, f); err != nil {
		log.Println(err)
	}

}

func handleMarkdownFile(filePath string, w io.Writer) error {

	md, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		w.Write([]byte("500"))
		return err
	}

	customHead := css +
		`<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
	<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3.0.1/es5/tex-mml-chtml.js"></script>
	<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.3.1/styles/atom-one-dark.min.css">
	<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.3.1/highlight.min.js"></script>
	<script>
	window.MathJax = {
		tex: {
			macros: {
			R : "\\mathbb{R}",
			map : "#1:#2\\rightarrow#3"
			}
		}
		};
	
		hljs.highlightAll();
		</script></head>
	`

	bytesReader := bytes.NewReader(md)
	bufReader := bufio.NewReader(bytesReader)
	title, _, err := bufReader.ReadLine()
	if nil != err {
		log.Fatal(err)
	}

	html := htmlPrefix(string(title), customHead)
	w.Write([]byte(html))

	rendered_md := markdown.ToHTML(md, nil, nil)
	w.Write(rendered_md)

	log.Println("done")
	return nil
}

func generateStaticBlog(name string) {

	postsPath := filepath.Join(".", "posts")
	blogPath := filepath.Join(".", "blog")

	postsDir, postsErr := os.Stat(postsPath)
	blogDir, blogErr := os.Stat(blogPath)

	if nil != postsErr || nil != blogErr {
		log.Fatal(`Missing "posts" and/or "blog" directories`)
	}

	if !postsDir.IsDir() || !blogDir.IsDir() {
		log.Fatal(`Missing "posts" and/or "blog" directories`)
	}

	files, err := ioutil.ReadDir(postsPath)
	if err != nil {
		log.Fatal(err)
	}

	var html = htmlPrefix(name, css)
	html += `<h2>` + name + `</h2><hr><ul>`

	for _, file := range files {
		var fileName = file.Name()
		if file.IsDir() {
			log.Fatal("make sure posts dir is flat")
		}

		filePath := filepath.Join(postsPath, fileName)
		log.Println("Handling " + filePath)

		fileNameNoExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		fileNameHtml := fileNameNoExt + ".html"

		htmlPostPath := filepath.Join(blogPath, fileNameHtml)
		f, err := os.Create(htmlPostPath)
		defer f.Close()

		w := bufio.NewWriter(f)

		err = handleMarkdownFile(filePath, w)
		if nil != err {
			log.Fatal(err)
		}

		w.Flush()

		t, err := times.Stat(filePath)
		if err != nil {
			log.Fatal(err.Error())
		}

		if !t.HasBirthTime() {
			log.Fatal("Could not get file's birthtime")
		}

		postDate := t.BirthTime().Format("2006-01-02")
		html += "<li><a href='" + fileNameHtml + "'>" + "<span style='color:black'>" + postDate + "</span>" + " " + fileNameNoExt + "</a></li>\n"
	}

	html += htmlPostfix
	buffer := []byte(html)

	indexPath := filepath.Join(blogPath, "index.html")
	err = ioutil.WriteFile(indexPath, buffer, 0644)

	if nil != err {
		log.Fatal(err)
	}
}

func runServer(portNumber int) error {

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	portNumberStr := strconv.Itoa(portNumber)
	localhostStr := "localhost"

	log.Println("[Here]")
	log.Println("Serving: " + path)
	log.Println("http://" + localhostStr + ":" + portNumberStr)

	http.HandleFunc("/", requestHandler)
	if err := http.ListenAndServe(localhostStr+":"+portNumberStr, nil); err != nil {
		return err
	}
	return nil
}

func main() {

	var portNumber int
	var makeBlog bool

	flag.IntVar(&portNumber, "port", 9898, "port number to listen on")
	flag.BoolVar(&makeBlog, "b", false, "generate static blog")
	flag.Parse()

	if makeBlog {
		log.Println("Generating static blog")
		generateStaticBlog("Blog")
		return
	}

	if err := runServer(portNumber); err != nil {
		log.Fatal(err)
	}

}
