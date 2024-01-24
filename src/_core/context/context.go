package context

import "go-trade/src/_core/config"

// TRADE
var BroadcastTrade = make(chan *TradeChanel, 1000)

var Config config.Config

func Init() {
	Config = config.GetConfig()
}
