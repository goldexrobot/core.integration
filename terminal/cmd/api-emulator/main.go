package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/activity"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/api"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/console"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/helpers"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		argPort    = flag.Uint("port", 8080, "Port to serve /ws")
		argConsole = flag.Bool("console", true, "Enable console interaction")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	// console
	var con *console.Console
	if *argConsole {
		c, clear, err := console.New()
		if err != nil {
			panic(err)
		}
		defer clear()
		con = c
	}

	// logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      con == nil,
		DisableColors:    con != nil,
		DisableTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	// console?
	if con != nil {
		log.SetOutput(con.Stderr())
		logger.SetOutput(io.Discard)
		logger.AddHook(helpers.LogrusConsoleHook{
			Consoler: con,
		})
	} else {
		logger.SetOutput(os.Stdout)
	}

	// api controller
	ctl := api.NewController(logger.WithField("api", "ctl"))

	// console?
	if con != nil {
		// "reading" event from console
		consoleReading := make(chan struct{}, 1)
		defer close(consoleReading)

		// serve console
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			con.Printf(`Welcome to Goldex terminal emulator. Try %v`, color.CyanString("help"))
			con.Run(ctx, &activity.Main{
				Logger: logger.WithField("activity", "main"),
				Ctl:    ctl,
			}, "", consoleReading)
		}()

		// wait console
		<-consoleReading
	}

	// api server: websocket + jsonrpc
	srv, err := api.NewServer(
		int(*argPort),
		apiv1.NewImpl(ctl, ctl, logger.WithField("api", "impl")),
		logger.WithField("api", "srv"),
	)
	if err != nil {
		fmt.Println("Failed to run server on port", *argPort)
		return
	}

	// serve http
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		srv.Serve(ctx)
	}()

	wg.Wait()
}
