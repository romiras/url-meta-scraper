package drivers

import (
	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/log"

	"github.com/streadway/amqp"
)

type AmqpConsumer struct {
	conn    *amqp.Connection
	topic   string
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
		done:   make(chan error),
		logger: logger,
	}, nil
}

func (pr *AmqpConsumer) Close() error {
	return pr.conn.Close()
}

func (pr *AmqpConsumer) Consume(handleFunc consumers.HandleFunc, logger log.Logger) error {
	var msgs <-chan amqp.Delivery
	// optionsMap := options.OptionsToMap(producerOptions)
	// metricName := options.GetOptionOrDefault(optionsMap, "metric_name", "task-producer").(string)

	ch, err := pr.getChannel()
	if err != nil {
		return err
	}
	defer func() {
		err = ch.Close()
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	// Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table)
	msgs, err = ch.Consume(
		pr.topic, // queue
		"",       // consumer tag
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return err
	}

	go func() {
		handleFunc(msgs /*pr.done,*/, logger)
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

/*
msgs, err := ch.Consume(
  q.Name, // queue
  "",     // consumer
  true,   // auto-ack
  false,  // exclusive
  false,  // no-local
  false,  // no-wait
  nil,    // args
)
failOnError(err, "Failed to register a consumer")

forever := make(chan bool)

go func() {
  for d := range msgs {
    log.Printf("Received a message: %s", d.Body)
  }
}()

log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
<-forever
*/

/*
   go func(con *amqp.Connection) {
       channel, _ := connection.Channel()
       defer channel.Close()
       durable, exclusive := false, false
       autoDelete, noWait := true, true
       q, _ := channel.QueueDeclare("test", durable, autoDelete, exclusive, noWait, nil)
       channel.QueueBind(q.Name, "#", "amq.topic", false, nil)
       autoAck, exclusive, noLocal, noWait := false, false, false, false
       messages, _ := channel.Consume(q.Name, "", autoAck, exclusive, noLocal, noWait, nil)
       multiAck := false
       for msg := range messages {
           fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
           msg.Ack(multiAck)
       }
   }(connection)
*/
