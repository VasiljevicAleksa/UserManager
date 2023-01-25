package rabbit

import (
	"usermanager/app/config"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type RMQ struct {
	PublishChannel chan uuid.UUID
}

// Create rabbitMQ connector and producer. Returns the channel where
// messages should be pushed in order to be published to rabbit queue.
func NewRMQ() *RMQ {
	publishChannel := make(chan uuid.UUID)
	initProducer(publishChannel)

	return &RMQ{
		PublishChannel: publishChannel,
	}
}

// Create rmq connection, open a channel and declare queue.
// If everything is done successfully, a new goroutine is started
// that listens to the channel where messages for publish are sent.
func initProducer(publishChannel chan uuid.UUID) {
	if config.EnvConfig.RabbitUrl == "" {
		log.Printf("Rabbit url is not defined")
	}

	conn, err := amqp.Dial(config.EnvConfig.RabbitUrl)
	if err != nil {
		log.Error().Err(err).Msg("cannot connect to rabbit")
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("cannot open rabbit channel")
		return
	}

	err = ch.ExchangeDeclare(
		config.EnvConfig.NotificationQueue, // name
		"fanout",                           // type
		true,                               // durable
		false,                              // auto-deleted
		false,                              // internal
		false,                              // no-wait
		nil,                                // arguments
	)

	if err != nil {
		log.Error().Err(err).Msg("cannot decalare rabbit exchange")
		return
	}

	log.Info().Msg("rmq ready to send messages")

	// start the listener process who will listen on publishChanell
	// for messages to be pushed to queue
	go listenForMessages(ch, publishChannel)
}

// listens on the channel for messages to be sent to the queue
func listenForMessages(ch *amqp.Channel, publishChannel chan uuid.UUID) {
	for userId := range publishChannel {
		err := ch.Publish(
			config.EnvConfig.NotificationQueue,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(userId.String()),
			},
		)

		if err != nil {
			log.Error().Err(err).Msg("cannot publish message")
			continue
		}

		log.Info().Msgf("notification for user %v successfully published", userId)
	}
}
