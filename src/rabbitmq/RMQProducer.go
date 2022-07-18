package rabbitmq

import (
	"encoding/json"
	"time"
	"user-ms/src/dto"

	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
)

type RMQProducer struct {
	ConnectionString string
}

func (r RMQProducer) StartRabbitMQ() (*amqp.Channel, error) {
	connectRabbitMQ, err := amqp.Dial(r.ConnectionString)

	if err != nil {
		return nil, err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()

	if err != nil {
		return nil, err
	}

	return channelRabbitMQ, err
}

func AddNotification(notification *dto.NotificationDTO, channel *amqp.Channel) {
	uuid, _ := uuid.NewV4()

	payload, _ := json.Marshal(notification)

	channel.Publish(
		"AddNotification-MS-exchange",    // exchange
		"AddNotification-MS-routing-key", // routing key
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.String(),
			Timestamp:    time.Now(),
			Body:         payload,
		})
}
