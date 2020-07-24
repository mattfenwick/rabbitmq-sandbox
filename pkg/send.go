package pkg

import (
	"fmt"
	errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	//"time"
)

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func Run() {
	//RunSend(myFirstChannel)
	//RunSend("another-channel")
	RunSend("channel-3")
}

func RunSend(queueName string) {
	args := os.Args
	fmt.Printf("%+v", args)

	rabbit, err := NewRabbit("guest", "guest", "localhost", 5672)
	doOrDie(errors.Wrapf(err, "Failed to connect to RabbitMQ"))

	log.Infof("connected to rabbitmq")

	ch, err := rabbit.Conn.Channel()
	doOrDie(errors.Wrapf(err, "Failed to open a channel"))
	defer ch.Close()

	log.Infof("opened a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	doOrDie(errors.Wrapf(err, "Failed to declare a queue"))

	log.Infof("declared a queue named %s", q.Name)

	for i := 0; i < 100000000000; i++ {
		message := fmt.Sprintf("Hello World! %d", i)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})
		doOrDie(errors.Wrapf(err, "Failed to publish a message"))

		log.Infof("published a message -- %s ", message)
		//time.Sleep(1 * time.Second)
	}
}
