package bootstrap

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

func Bootstrap() {
	fmt.Println(banner)

	cfg, err := config.Initialize()
	if err != nil {
		logrus.Fatalf("Could not parse configuration file (%s): %v", config.ConfigLocation, err)
		return
	}

	records, err := rec.Initialize(cfg)
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

	err = speedtest.StartSpeedTestService(records, cfg, ctx, wg)
	if err != nil {
		logrus.Fatalf("Could not start speedtest service: %v", err)
		cancel()
		goto close
	}

	err = api.StartHTTPServer(cfg, records, ctx)
	if err != nil {
		logrus.Fatalf("Error occurred in HTTP server: %v", err)
		cancel()
		goto close
	}

close:

	wg.Wait()
}
