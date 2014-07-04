package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/featen/ags/service"
	"github.com/featen/ags/service/config"
	log "github.com/featen/utils/log"
)

const (
	servicePrefix = "/service"
	indexFile     = "index.html"
)

var (
	staticDir = http.Dir("./webapp")
)

func staticFileshandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		return
	}
	file := req.URL.Path

	if strings.HasPrefix(file, servicePrefix) {
		log.Debug("request service")
		return
	}

	if file != "" && file[0] != '/' {
		return
	}

	f, err := staticDir.Open(file)
	if err != nil {
		http.NotFound(res, req)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(req.URL.Path, "/") {
			http.Redirect(res, req, req.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, indexFile)
		f, err = staticDir.Open(file)
		if err != nil {
			http.NotFound(res, req)
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return
		}
	}

	log.Debug("[Static] Serving " + file)

	http.ServeContent(res, req, file, fi.ModTime(), f)
}

func clean() {
	log.Fatal("db & log clean, exiting...")
	log.Close()
	os.Exit(0)
}

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sc
		log.Info("signal recieved: %v", s)
		clean()
	}()

	service.RegService()

	http.HandleFunc("/", staticFileshandler)

	var bind string
	if len(os.Getenv("HOST")) == 0 || len(os.Getenv("PORT")) == 0 {
		bind = fmt.Sprintf(":%s", config.GetValue("ServicePort"))
	} else {
		bind = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	}

	log.Info("start listening on %s", bind)
	fmt.Println("ags server started on %s ", bind)

	http.ListenAndServe(bind, nil)
}
