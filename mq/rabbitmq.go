package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"rscheduler/config"
	"rscheduler/global"
	"time"
)

type rabbitMQ struct {
	conn *amqp.Connection
}

var RabbitMQ *rabbitMQ

func Init() {
	RabbitMQ = rabbitMQInit()
}

func rabbitMQInit() *rabbitMQ {
	conn, err := amqp.Dial(config.Config.RabbitMQ.URL)
	if err != nil {
		panic("get rabbitmq conn failed, err: " + err.Error())
	}
	return &rabbitMQ{conn: conn}
}

// Get it will get a task from rabbitmq
// When the queue is empty, it will block and try again after 1 second
func (r *rabbitMQ) Get() ([]byte, error) {
	channel, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			global.Logger.Error("Close rabbitmq channel failed, err: " + err.Error())
		}
	}(channel)

	// Get a task from rabbitmq
	cfg := config.Config.RabbitMQ
	for {
		d, ok, err := channel.Get(cfg.TaskQueueName, false)
		if err != nil {
			return nil, err
		}

		if ok {
			global.Logger.Info("Get a new task from rabbitmq")
			// Ack the task
			err = d.Ack(false)
			if err != nil {
				global.Logger.Error("Ack task failed, err: " + err.Error())
			}
			return d.Body, nil
		}

		// if not ok, try again after 1 second
		time.Sleep(time.Second)
	}
}
