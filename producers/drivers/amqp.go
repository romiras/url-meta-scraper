package drivers

import (
	"io/ioutil"
	"time"

	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/options"
	"github.com/romiras/url-meta-scraper/producers"

	"github.com/streadway/amqp"
)

type AmqpProducer struct {
	conn    *amqp.Connection
	topic   string
	channel *amqp.Channel
	logger  log.Logger
}

func NewAmqpProducer(amqpURI, topic string, logger log.Logger) (producers.TaskProducer, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		// logger.Fatal("Unable connect to RabbitMQ", "error", err.Error())
		return nil, err
	}

	logger.Info("Connection to RabbitMQ successful")
	return &AmqpProducer{
		conn:   conn,
		topic:  topic,
		logger: logger,
	}, nil
}

func (pr *AmqpProducer) Close() error {
	return pr.conn.Close()
}

func (pr *AmqpProducer) Produce(task *producers.Task, producerOptions ...options.Option) (string, error) {
	// optionsMap := options.OptionsToMap(producerOptions)
	// metricName := options.GetOptionOrDefault(optionsMap, "metric_name", "task-producer").(string)
	var b []byte

	ch, err := pr.getChannel()
	if err != nil {
		return "", err
	}
	defer ch.Close()

	b, err = ioutil.ReadAll(task.Body)
	if err != nil {
		return "", err
	}

	pub := amqp.Publishing{
		Body:         b,
		Headers:      amqp.Table{},
		ContentType:  "application/json",
		DeliveryMode: 2, // Persistent
		Timestamp:    time.Now().UTC(),
	}

	// Publish(exchange, key string, mandatory, immediate bool, msg Publishing)
	err = ch.Publish("", pr.topic, false, false, pub)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (pr *AmqpProducer) getChannel() (*amqp.Channel, error) {
	var err error
	pr.channel, err = pr.conn.Channel()
	if err != nil {
		return nil, err
	}

	return pr.channel, nil
}
