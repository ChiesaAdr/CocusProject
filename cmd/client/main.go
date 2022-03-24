package main

import (
	"log"

	"cocus/internal/config"
	"cocus/internal/service"
	"cocus/internal/service/client"
	"cocus/internal/ui"
)

func main() {

	//Set default to log show(time pattern)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Loading .env to get local variables
	config.LoadDotEnv()

	//Create a service of CLIENT type
	s := service.NewService(service.CLIENT).(*client.Client)
	//TODO: Add the possibility to get this address through the executable argument
	s.Run(config.GetAddressServer())

	//Start the User Interface
	//TODO: Create a modo to select other types(Websocket, basic terminal)
	ui.StartUi(s)
}
