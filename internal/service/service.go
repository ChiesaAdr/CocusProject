package service

import (
	"cocus/internal/service/client"
	"cocus/internal/service/server"
)

const (
	SERVER int = iota
	CLIENT
)

type Service interface {
	Run(addr string)
}

func NewService(stype int) Service {
	var ser Service

	if stype == SERVER {
		ser = server.NewServer()
	} else if stype == CLIENT {
		ser = client.NewClient()
	}
	return ser
}
