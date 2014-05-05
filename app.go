package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/featen/ags/service"
	"github.com/featen/ags/service/config"
	log "github.com/featen/utils/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	conf          map[string]string
	staticFilesWs *restful.WebService
)

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

	log.Info("start listening on localhost:%s", config.GetValue("ServicePort"))
	fmt.Println("ags server started on port ", config.GetValue("ServicePort"))
	http.ListenAndServe(":"+config.GetValue("ServicePort"), nil)
}
