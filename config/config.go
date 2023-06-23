package config

var Config = &configModel{}

type configModel struct {
	RabbitMQ struct {
		URL string `json:"url"`
	} `json:"rabbitmq"`
}
