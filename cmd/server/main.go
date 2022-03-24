package main

import (
	"cocus/internal/config"
	"cocus/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Channel to wait a signal and finish the server
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	//Set default to log show(time pattern)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Loading .env to get local variables
	config.LoadDotEnv()

	//Create a service of SERVER type
	s := service.NewService(service.SERVER)

	s.Run(config.GetAddressServer())

	log.Println("Server is up!")
	<-done
	log.Println("(⌐■_■) Server is out!")
}
