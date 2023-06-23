package config

var Config = &configModel{}

type configModel struct {
	RabbitMQ struct {
		URL           string `json:"url"`
		TaskQueueName string `json:"taskQueueName"`
	} `json:"rabbitmq"`
	TaskLimit struct {
		MaxCPU int `json:"maxCPU"`
		MaxMem int `json:"maxMem"`
	} `json:"taskLimit"`
}
