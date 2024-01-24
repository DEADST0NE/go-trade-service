package rabbit

import (
	"encoding/json"
	"time"

	"go-trade/src/_core/context"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func TradePublisher() {
	log.Info("START TRADE PUBLISHER")

	for {
		conn, err := Connect()
		if err != nil {
			log.Errorf("Failed to connect to RabbitMQ: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Errorf("Failed to open a channel: %s", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		queueName := "Trade"
		exchangeName := "TradeExchange"

		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Errorf("Failed to declare a queue: %s", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		err = ch.ExchangeDeclare(
			exchangeName,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Errorf("Failed to declare an exchange: %s", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		ch.QueueBind(q.Name, q.Name+"."+"*", exchangeName, false, nil)

		for {
			trade := <-context.BroadcastTrade

			tradeBytes, err := json.Marshal(trade)
			if err != nil {
				log.Errorf("Failed to serialize trade: %s", err)
				continue
			}

			routingKey := q.Name + "." + trade.Symbol

			err = ch.Publish(
				exchangeName,
				routingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        tradeBytes,
				})

			if err != nil {
				log.Errorf("Failed to publish a message: %s", err)
				break
			}
		}

		ch.Close()
		conn.Close()
		time.Sleep(5 * time.Second)
	}
}
