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
	url       = "amqp://lxs:root@192.168.1.4:5672/"
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
			InputFile:        "http://127.0.0.1:9000/hiplot-copilot/data.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=FNNRPZ94ZSRCOHKILP7W%2F20240324%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240324T153323Z&X-Amz-Expires=604800&X-Amz-Security-Token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJGTk5SUFo5NFpTUkNPSEtJTFA3VyIsImV4cCI6MTcxMTI5ODYzMywicGFyZW50IjoibWluaW9hZG1pbiJ9.hV4F6OJMsia5C4wltzmMzIQp-jaDvrEY2VfIoEFR4-km1eDpoTFc8HkDa0v5RPvMz6LfroAty2UMbDcwo839VQ&X-Amz-SignedHeaders=host&versionId=null&X-Amz-Signature=08c6c9049547bb2fe0d9e4189715d7d9a9a1fc0cbab864ba6d0f4540e9a405c3",
			ConfFile:         "http://127.0.0.1:9000/hiplot-copilot/data.json?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=FNNRPZ94ZSRCOHKILP7W%2F20240324%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240324T153846Z&X-Amz-Expires=604800&X-Amz-Security-Token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJGTk5SUFo5NFpTUkNPSEtJTFA3VyIsImV4cCI6MTcxMTI5ODYzMywicGFyZW50IjoibWluaW9hZG1pbiJ9.hV4F6OJMsia5C4wltzmMzIQp-jaDvrEY2VfIoEFR4-km1eDpoTFc8HkDa0v5RPvMz6LfroAty2UMbDcwo839VQ&X-Amz-SignedHeaders=host&versionId=null&X-Amz-Signature=12d66f481c6f85b59b9dc428470322f9ba2ce0c1076dd6bc67a05afa3b267129",
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
