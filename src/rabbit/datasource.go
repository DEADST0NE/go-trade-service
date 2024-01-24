package rabbit

import (
	"go-trade/src/_core/context"

	"github.com/streadway/amqp"
)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(context.Config.Rabbit.Url)

	return conn, err
}
