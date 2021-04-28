package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/zhangjie2012/yinyang/pkg/api"
	"github.com/zhangjie2012/yinyang/pkg/calendar"
)

var (
	Version     string = ""
	Branch      string = ""
	BuildTime   string = ""
	rawDataPath        = "/data/data/rawdata"
	host               = "localhost"
	port               = 9123
)

func init() {
	flag.StringVar(&rawDataPath, "rawdata", rawDataPath, "the yin-yang relate rawdata file path")
	flag.StringVar(&host, "host", host, "the server host")
	flag.IntVar(&port, "port", port, "the server port")

	callerPrettyfier := func(frame *runtime.Frame) (function string, file string) {
		ss := strings.Split(frame.Function, ".")
		function = ss[len(ss)-1]
		file = fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		return function, file
	}
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: callerPrettyfier,
		TimestampFormat:  "2006-01-02 15:04:05.000",
		FullTimestamp:    true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infof("build info: branch=%s, version=%s, buildtime=%s", Branch, Version, BuildTime)

	flag.Parse()

	calendar.ParseRawData(rawDataPath)
	log.Infof("parse raw data success")

	wg := sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start server
	server := api.NewServer(host, port)
	server.Start(ctx, &wg)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infof("receiver shutdown signal")
}
