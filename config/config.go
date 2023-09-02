package config

var Config = &configModel{}

type configModel struct {
	RabbitMQ struct {
		URL           string
		TaskQueueName string
	}
	TaskLimit struct {
		MaxCPU              int
		MaxMem              int
		MaxBusyProcessor    int
		MaxIdleProcessor    int
		MaxBusyProcessorMem int
		MaxIdleProcessorMem int
		TaskTimeout         int
	}
}
