package worker

import (
	"stacktack-go/src/worker"
	"testing"
	"time"
)

func TestNovaConsumerRun(t *testing.T) {
	config := worker.NewConfigByJson("E:\\gopath\\src\\stacktack-go\\test\\worker\\config.json")

	novaConsumer, err := worker.NewNovaConsumer(config)
	if err != nil {
		t.Error(err)
		return
	}
	novaConsumer.Run(worker.NovaNotificationInfoHandler)

}

func TestNovaConsumerShutDown(t *testing.T) {
	config := worker.NewConfigByJson("E:\\gopath\\src\\stacktack-go\\test\\worker\\config.json")

	novaConsumer, err := worker.NewNovaConsumer(config)
	if err != nil {
		t.Error(err)
		return
	}
	go novaConsumer.Run(worker.NovaNotificationInfoHandler)
	time.Sleep(10 * time.Second) //make sure shutdown should be called after Run
	novaConsumer.Shutdown(novaConsumer.ConsumerTag)
}
