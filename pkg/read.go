package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	myFirstChannel = "hello"
)

func RunRead() {
	ReadFromChannel(myFirstChannel, make(chan struct{}))
}

type Rabbit struct {
	Username string
	Password string
	Address  string
	Port     int
	Conn     *amqp.Connection
}

func NewRabbit(username string, password string, address string, port int) (*Rabbit, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, address, port))
	if err != nil {
		return nil, err
	}
	return &Rabbit{
		Username: username,
		Password: password,
		Address:  address,
		Port:     port,
		Conn:     conn,
	}, nil
}

func (r *Rabbit) Close() error {
	if r.Conn == nil {
		return errors.Errorf("connection already closed")
	}
	r.Conn.Close()
	r.Conn = nil
	return nil
}

func ReadFromChannel(queueName string, stop <-chan struct{}) {
	rabbit, err := NewRabbit("guest", "guest", "localhost", 5672)
	doOrDie(errors.Wrapf(err, "Failed to connect to RabbitMQ"))

	ch, err := rabbit.Conn.Channel()
	doOrDie(errors.Wrapf(err, "Failed to open a channel"))
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	doOrDie(errors.Wrapf(err, "Failed to declare a queue"))

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	doOrDie(errors.Wrapf(err, "Failed to register a consumer"))

	go func() {
		for d := range msgs {
			log.Infof("Received a message: %s", d.Body)
		}
	}()

	log.Infof(" [*] Waiting for messages. To exit press CTRL+C")
	<-stop
}
