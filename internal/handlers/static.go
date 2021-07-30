package handlers

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var httpFileSystem http.FileSystem
var fsHandler http.Handler

var fileSystem fs.FS

var reLabel = regexp.MustCompile(`\{\{\s*label\s*\}\}`)
var label_element = `<div id="label"></div>`

var reBodyClose = regexp.MustCompile(`</\s*body\s*>`)
var bodyClose = `<script src="/browser-view.js"></script></body>`

var reHeadClose = regexp.MustCompile(`</\s*head\s*>`)
var headClose = `<meta charset="utf-8">
<meta name="viewport" content="width=device-width, , initial-scale=1.0, maximum-scale=1.0, user-scalable=0">
<!-- favicon and webmanifest -->
<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="manifest" href="/site.webmanifest">
<title>Label Browser View</title></head>`

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fsHandler.ServeHTTP(w, r)
	log.Printf("[web] GET 200 %s", r.URL)
}

func StaticHandlerBrowserView(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	layout := qs.Get("layout")

	var str string

	if layout != "" && layout != "default" {
		if !strings.HasSuffix(layout, ".html") {
			layout = layout + ".html"
		}
		fp := filepath.Join(*Cfg.LabelDirectory, "layouts", layout)
		if b, e := os.ReadFile(fp); e == nil {
			str = string(b)
			log.Printf("[web] GET 200 %s", r.URL)
		}
	}

	if str == "" {
		b, _ := fs.ReadFile(fileSystem, filepath.Join("browser-view.html"))
		str = string(b)
		log.Printf("[web] GET 200 %s layout not found - used default", r.URL)
	}

	str = reLabel.ReplaceAllString(str, label_element)
	str = reHeadClose.ReplaceAllString(str, headClose)
	str = reBodyClose.ReplaceAllString(str, bodyClose)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(str))
}

func GetFileSystem(embeddedFiles fs.FS, cfgDev *bool) {
	useEmbedded := true
	if cfgDev != nil && *cfgDev {
		useEmbedded = false
	}

	if !useEmbedded {
		if _, err := os.Stat("static"); err == nil {
			log.Print("[load] Using static files directly for web content")
			httpFileSystem = http.FS(os.DirFS("static"))
			fileSystem = os.DirFS("static")
			fsHandler = http.FileServer(httpFileSystem)
			return
		}
	}

	log.Print("[load] Using embedded files for web content")

	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}

	fileSystem = fsys
	httpFileSystem = http.FS(fsys)
	fsHandler = http.FileServer(httpFileSystem)
}
