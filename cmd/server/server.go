package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitbucket.org/latonaio/authenticator/configs"
	"bitbucket.org/latonaio/authenticator/pkg/db"
	"bitbucket.org/latonaio/authenticator/pkg/server"
)

//var (
//	Version  string
//	Revision string
//)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfgs, err := configs.New()
	if err != nil {
		panic(err)
	}
	echoServer := server.New(ctx, cfgs)
	err = db.NewDBConPool(ctx, cfgs)
	if err != nil {
		panic(err)
	}

	errC := make(chan error)
	// Start server
	go echoServer.Start(errC)

	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGTERM, os.Interrupt)

	select {
	case err := <-errC:
		panic(err)
	case <-quitC:
		if err := echoServer.Shutdown(ctx); err != nil {
			errC <- err
		}
		cancel()
		time.Sleep(1 * time.Second)
	}
}
