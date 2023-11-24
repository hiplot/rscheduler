package mq

import (
	"context"
	"encoding/json"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

type HiPlotTask struct {
	InputFile        string `json:"inputFile"`
	ConfFile         string `json:"confFile"`
	OutputFilePrefix string `json:"outputFilePrefix"`
	Tool             string `json:"tool"`
	Module           string `json:"module"`
	ID               string `json:"ID"`
	Name             string `json:"Name"` // 为了适配多种任务，common为通用，从plumber移植来的，执行commonInit.R
}

const (
	queueName = "common"
	url       = "amqp://lxs:root@172.30.0.1:5672/"
	COUNT     = 1
)

func TestPushMsg(t *testing.T) {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic("get rabbitmq conn failed, err: " + err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		t.Error("get channel failed, err: " + err.Error())
		return
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		t.Error("declare queue failed, err: " + err.Error())
		return
	}

	for i := 0; i < COUNT; i++ {
		id, _ := gonanoid.New()
		task := HiPlotTask{
			InputFile:        "/home/lxs/code/ospp/user/input/Ps3qMcfKcGp_uruEC1x-A/data.txt",
			ConfFile:         "/home/lxs/code/ospp/user/input/Ps3qMcfKcGp_uruEC1x-A/data.json",
			OutputFilePrefix: "/home/lxs/code/ospp/user/output/" + id,
			Tool:             "interval-area-chart",
			Module:           "basic",
			ID:               id,
			Name:             "common",
		}
		body, err := json.Marshal(&task)
		if err != nil {
			t.Error("json marshal failed, err: " + err.Error())
			return
		}
		err = ch.PublishWithContext(context.Background(),
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})
		if err != nil {
			t.Error("publish msg failed, err: " + err.Error())
			return
		}
	}
}
