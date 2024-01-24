package broker

import (
	"encoding/json"
	"go-trade/src/_core/context"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

func ssConnectToServer() (*websocket.Conn, error) {
	url := url.URL{
		Scheme: context.Config.Broker.Ss.Scheme,
		Path:   context.Config.Broker.Ss.Path,
		Host:   context.Config.Broker.Ss.Host,
	}
	log.Info("СONNECTING TO SS WS URL: ", url.String())

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil {
		return nil, err
	}

	subMessage := SsSubMessage{
		Method: "subscribe",
		Ch:     "trades",
		Params: SsParams{
			Symbols: context.Config.Symbols,
		},
	}

	messageBytes, err := json.Marshal(subMessage)

	if err != nil {
		return nil, err
	}

	c.WriteMessage(websocket.TextMessage, messageBytes)

	if err != nil {
		return nil, err
	}

	return c, err
}

func ssListenAndServe(c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()

		if err != nil {
			log.Error("Error read SS ws:", err)
			c.Close()
			return
		}

		var message SsMsgTrades
		err = json.Unmarshal(msg, &message)

		if err != nil {
			log.Error("Error parse message SS ws:", err)
		} else {

			if message.Update != nil {
				for symbol, trades := range *message.Update {
					for _, trade := range trades {
						msg := context.TradeChanel{
							Timestamp: trade.T,
							Id:        trade.I,
							Price:     trade.P,
							Quantity:  trade.Q,
							Side:      trade.S,
							Symbol:    symbol,
						}

						context.BroadcastTrade <- &msg
					}
				}
			}
		}
	}
}

func SSListener() {
	const maxAttempts = 5
	attempts := 0

	for {
		c, err := ssConnectToServer()
		if err != nil {
			log.Error("Error connect SS ws:", err)
		} else {
			ssListenAndServe(c)
			attempts = 0
		}

		if attempts++; attempts >= maxAttempts {
			log.Error("Maximum number of connection SS ws attempts exceeded:", err)
			os.Exit(1)
			break
		}

		log.Info("Attempting to reconnect after 5 seconds...:")
		time.Sleep(5 * time.Second)
	}
}
