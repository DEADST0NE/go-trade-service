package rabbit

import (
	"go-trade/src/_core/context"

	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func TradePublisher() {
	log.Info("START TRADE PUBLISHER")
	queueName := "Trade"
	exchangeName := "TradeExchange"

	conn, err := Connect()

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	ch.QueueBind(queueName, q.Name+"."+"*", exchangeName, false, nil)

	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	for {
		trade := <-context.BroadcastTrade

		tradeBytes, err := json.Marshal(trade)
		if err != nil {
			log.Fatalf("Failed to serialize trade: %s", err)
		}

		routingKey := q.Name + "." + trade.Symbol

		err = ch.Publish(
			exchangeName, // exchange
			routingKey,   // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        tradeBytes,
			})

		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}

	}
}
