package main

import (
	"go-trade/src/_core/context"
	"go-trade/src/broker"
	"go-trade/src/rabbit"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("START PROJECT")

	// CONTEXT
	context.Init()

	// PUBLISHER
	go rabbit.TradePublisher()

	// DATASOURCE
	broker.SSListener()
}
