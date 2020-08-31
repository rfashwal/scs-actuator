package mq

import (
	"encoding/json"
	"fmt"

	"github.com/rfashwal/scs-actuator/internal"
	"github.com/rfashwal/scs-actuator/internal/dto"
	"github.com/rfashwal/scs-utilities/config"
	"github.com/rfashwal/scs-utilities/rabbit"
)

func TemperatureObserver(s internal.Service, manager rabbit.MQManager, conf config.Manager) error {
	observer, err := manager.InitObserver()
	if err != nil {
		return err
	}
	defer manager.CloseConnection()
	defer observer.Channel.Close()

	err = observer.DeclareTopicExchange(conf.ActuatorTopic())
	if err != nil {
		return err
	}
	err = observer.BindQueue(observer.Queue, conf.ReadingsRoutingKey()+".#", conf.ActuatorTopic())

	if err != nil {
		return err
	}

	deliveries := observer.Observe()

	for msg := range deliveries {
		aculatorMessage := dto.AculatorMessage{}
		err := json.Unmarshal(msg.Body, &aculatorMessage)
		if err != nil {
			fmt.Printf("could not unmarshal expected aculator message, %s\n", err.Error())
			continue
		}
		fmt.Printf("message is delivered %v\n", aculatorMessage)
		// Do Something with this value, update Valve
	}
	return nil
}
