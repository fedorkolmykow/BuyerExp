package main

import (
	"context"
	"github.com/fedorkolmykow/avitoexp/pkg/parser"
	"github.com/fedorkolmykow/avitoexp/pkg/smtp"
	"github.com/fedorkolmykow/avitoexp/pkg/worker"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fedorkolmykow/avitoexp/pkg/httpServer"
	"github.com/fedorkolmykow/avitoexp/pkg/dbClient"
	"github.com/fedorkolmykow/avitoexp/pkg/service"

	log "github.com/sirupsen/logrus"
)



func main() {
	log.SetFormatter(&log.JSONFormatter{})
	switch os.Getenv("LOG_LEVEL"){
		case "TRACE": log.SetLevel(log.TraceLevel)
		case "WARN": log.SetLevel(log.WarnLevel)
		case "FATAL": log.SetLevel(log.FatalLevel)
		default: log.SetLevel(log.FatalLevel)
	}
	err := os.Mkdir("logs", 0777)
	if err != nil {
		log.Warn(err)
	}
	file, err := os.OpenFile("logs/auto.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
	    log.Warn("Failed to log to file, using default stderr")
	}

	smtpCl := smtp.NewSMTP()
    dbCon := dbClient.NewDbClient()
    par := parser.NewParser()

    wrk := worker.NewWorker(dbCon, smtpCl, par)
    swc := service.NewService(dbCon, smtpCl)
	router := httpServer.NewHTTPServer(swc)
	srv := &http.Server{
		Addr:    os.Getenv("HTTP_PORT"),
		Handler: router,
	}

	go func() {
		log.Trace("starting HTTP server at", os.Getenv("HTTP_PORT"))
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed{
			log.Fatal(err)
		}
	}()

	go func() {
		t := os.Getenv("TIME_BETWEEN_PARSING")
		d, err := time.ParseDuration(t)
		if err != nil{
			log.Warn(err)
			d = time.Hour
		}
		for {
			err := wrk.SendChanges()
			if err != nil {
				log.Warn(err)
			}
			time.Sleep(d)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	wait, err := strconv.Atoi(os.Getenv("TIME_TO_SHUTDOWN"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(wait)*time.Second)
	defer func(){
		e := dbCon.Shutdown()
		if e != nil{
			log.Warn(e)
		}
		cancel()
	}()
	err = srv.Shutdown(ctx)
	if err != nil{
		log.Fatalf("Graceful Server Shutdown Failed:%+v", err)
	}
	log.Trace("Server Was Gracefully Stopped")
}