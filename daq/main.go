package main

import (
	"context"
	"daq/internal/api"
	"daq/internal/config"
	"daq/internal/rec"
	"daq/internal/speedtest"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

const banner = `
    ____        __        
   / __ \__  __/ /_______ 
  / /_/ / / / / / ___/ _ \
 / ____/ /_/ / (__  )  __/
/_/    \__,_/_/____/\___/ 
                          `

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	fmt.Println(banner)

	cfg := config.ParseConfig()
	records, err := rec.Initialize()
	if err != nil {
		logrus.Fatalf("Could not initialize records: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-closeSignal
		cancel()
	}()

	wg := &sync.WaitGroup{}

	speedtest.StartSpeedTestService(records, cfg, ctx, wg)

	err = api.StartHTTPServer(cfg, records, ctx)
	if err != nil {
		logrus.Fatalf("Error occurred in HTTP server: %v", err)
		return
	}

	wg.Wait()

	logrus.Info("Goodbye!")
}
