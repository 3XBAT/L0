package nats

import (
	"L0"
	"L0/pkg/repository"
	"L0/pkg/service"
	"encoding/json"

	"github.com/nats-io/stan.go"

	"github.com/sirupsen/logrus"
)

func NewSubscribeToChannel(clusterId, clientId, channelName string, repository *repository.Repository, service *service.Service) (stan.Conn, error) {

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		logrus.Errorf("Error while connecting to NATS-streaming: %s", err.Error())
	}

	msgHandler := func(msg *stan.Msg) {
		var order L0.Order

		if err := json.Unmarshal(msg.Data, &order); err != nil {
			logrus.Errorf("Error while unparsing JSON: %s", err.Error())
			return
		}

		if err := repository.SaveOrder(order); err != nil {
			logrus.Errorf("Error while pasting in database : %s", err.Error())
			return
		}

		service.Cache.AddOrder(order.OrderUID, order)

	}
	_, err = sc.Subscribe(channelName, msgHandler)

	if err != nil {
		logrus.Errorf("Error while subscribing to channel: %s", err.Error())
		return nil, err
	}
	logrus.Printf("Succesful connect")
	return sc, nil
}
