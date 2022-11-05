package main

import (
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/synthia-telemed/push-notification-consumer/pkg/config"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func failOnError(logger *zap.SugaredLogger, err error, msg string) {
	if err != nil {
		logger.Fatalw(msg, err)
	}
}

const (
	QueueName = "push-notification"
)

func main() {
	z, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Failed to creat zap logger")
	}
	logger := z.Sugar()

	cfg, err := config.Load()
	failOnError(logger, err, "Failed to parse env")

	conn, err := amqp.Dial(cfg.RabbitMQ.GetURL())
	failOnError(logger, err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(logger, err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(logger, err, "Failed to declare a queue")
	consumerName := uuid.NewString()
	msgs, err := ch.Consume(
		q.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(logger, err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			logger.Infow("Received msg", "body", d.Body)
			time.Sleep(time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()
	logger.Infow("Started listening for the message", "queue", q.Name, "consumer_name", consumerName)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down consumer...")
	failOnError(logger, ch.Cancel(consumerName, false), "Failed to cancel the channel")
	failOnError(logger, ch.Close(), "Failed to close the channel")
	failOnError(logger, conn.Close(), "Failed to close the connection")
	logger.Info("Consumer exiting")
}
