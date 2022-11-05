package main

import (
	"context"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/synthia-telemed/push-notification-consumer/pkg/config"
	"github.com/synthia-telemed/push-notification-consumer/pkg/consumer"
	"github.com/synthia-telemed/push-notification-consumer/pkg/notification"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func assertError(logger *zap.SugaredLogger, err error, isFatal bool, msg string) bool {
	if err == nil {
		return false
	}
	sentry.CaptureException(err)
	if isFatal {
		sentry.Flush(time.Second * 2)
		logger.Fatalw(msg, "error", err.Error())
		return true
	}
	logger.Errorw(msg, "error", err.Error())
	return true
}

func main() {
	z, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Failed to creat zap logger")
	}
	logger := z.Sugar()
	cfg, err := config.Load()
	assertError(logger, err, true, "Failed to parse env")
	err = sentry.Init(sentry.ClientOptions{TracesSampleRate: 1.0, Dsn: cfg.SentryDSN})
	assertError(logger, err, false, "Failed to init sentry")
	notificationClient, err := notification.NewFirebaseNotificationClient(context.Background(), &cfg.Notification)
	assertError(logger, err, true, "Failed to init notification client")

	conn, err := amqp.Dial(cfg.RabbitMQ.GetURL())
	assertError(logger, err, true, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	assertError(logger, err, true, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	assertError(logger, err, true, "Failed to declare a queue")
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
	assertError(logger, err, true, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var msg consumer.PushNotificationMessageBody
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logger.Infow("parse failed", "msg", err.Error())
				d.Nack(false, false)
				continue
			}
			if !msg.IsValid() {
				logger.Info("invalid")
				d.Nack(false, false)
				continue
			}
			params := notification.SendParams{
				Token: msg.Token,
				Title: msg.Title,
				Body:  msg.Body,
			}
			err := notificationClient.Send(context.Background(), params, msg.Data)
			if assertError(logger, err, false, "notificationClient.Send error") {
				d.Nack(false, false)
				continue
			}
			d.Ack(false)
		}
	}()
	logger.Infow("Started listening for the message", "queue", q.Name, "consumer_name", consumerName)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down consumer...")
	assertError(logger, ch.Cancel(consumerName, false), true, "Failed to cancel the channel")
	assertError(logger, ch.Close(), true, "Failed to close the channel")
	assertError(logger, conn.Close(), true, "Failed to close the connection")
	logger.Info("Consumer exiting")
}
