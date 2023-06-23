package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"rscheduler/config"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

func RabbitMQInit() *RabbitMQ {
	conn, err := amqp.Dial(config.Config.RabbitMQ.URL)
	if err != nil {
		panic("get rabbitmq conn failed, err: " + err.Error())
	}
	return &RabbitMQ{conn: conn}
}
