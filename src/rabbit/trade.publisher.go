package rabbit

import (
	"go-trade/src/_core/context"

	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func getTradeExchange(symbol string) string {
	return "Trade.exchange" + "." + symbol
}

func TradePublisher() {
	log.Info("START TRADE PUBLISHER")
	queueName := "Trade"

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

	for _, symbol := range context.Config.Symbols {
		err = ch.ExchangeDeclare(
			getTradeExchange(symbol), // name
			"topic",                  // type
			true,                     // durable
			false,                    // auto-deleted
			false,                    // internal
			false,                    // no-wait
			nil,                      // arguments
		)

		if err != nil {
			log.Fatalf("Failed to declare an exchange: %s", err)
		}
	}

	for {
		trade := <-context.BroadcastTrade

		tradeBytes, err := json.Marshal(trade)
		if err != nil {
			log.Fatalf("Failed to serialize trade: %s", err)
		}

		err = ch.Publish(
			getTradeExchange(trade.Symbol), // exchange
			q.Name,                         // routing key
			false,                          // mandatory
			false,                          // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        tradeBytes,
			})

		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}

	}
}

func SubscribeToTradeMessages(exchangeName string) {
	log.Info("START TRADE SUBSCRIBE")
	queueName := "Trade"

	conn, err := Connect()

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	defer ch.Close()

	// Объявление очереди
	q, err := ch.QueueDeclare(
		queueName, // имя очереди
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Привязка очереди к exchange
	err = ch.QueueBind(
		q.Name,       // имя очереди
		"#",          // routing key (подписка на все сообщения)
		exchangeName, // имя exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // очередь
		"",     // потребитель
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // аргументы
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
}
