package drivers

import (
	"fmt"

	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/log"

	"github.com/streadway/amqp"
)

type AmqpConsumer struct {
	conn    *amqp.Connection
	topic   string
	tag     string
	channel *amqp.Channel
	done    chan error
	logger  log.Logger
}

func NewAmqpConsumer(amqpURI, topic string, logger log.Logger) (consumers.IConsumer, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		// logger.Fatal("Unable connect to RabbitMQ", "error", err.Error())
		return nil, err
	}

	logger.Info("Connection to RabbitMQ successful")
	return &AmqpConsumer{
		conn:   conn,
		topic:  topic,
		tag:    "consumer-1",
		done:   make(chan error),
		logger: logger,
	}, nil
}

func (pr *AmqpConsumer) Close() error {
	return pr.conn.Close()
}

func (pr *AmqpConsumer) Consume(handleFunc consumers.HandleFunc, logger log.Logger) error {
	var err error
	// optionsMap := options.OptionsToMap(producerOptions)
	// metricName := options.GetOptionOrDefault(optionsMap, "metric_name", "task-producer").(string)

	logger.Info("got Connection, getting Channel")
	pr.channel, err = pr.getChannel()
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Queue bound to Exchange, starting Consume (consumer tag '%s')", pr.tag))
	deliveries, err := pr.channel.Consume(
		pr.topic, // name
		pr.tag,   // consumerTag,
		false,    // noAck
		false,    // exclusive
		false,    // noLocal
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("Queue Consume: %s", err))
	}

	go func() {
		handleFunc(deliveries, logger)
		pr.done <- nil
	}()

	return nil
}

func (pr *AmqpConsumer) getChannel() (*amqp.Channel, error) {
	// if pr.channel != nil {
	// 	return pr.channel, nil
	// }

	var err error
	pr.channel, err = pr.conn.Channel()
	if err != nil {
		return nil, err
	}

	return pr.channel, nil
}
