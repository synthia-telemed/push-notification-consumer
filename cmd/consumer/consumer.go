package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
	"github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
	"github.com/synthia-telemed/push-notification-consumer/pkg/notification"
	"go.uber.org/zap"
)

type Consumer interface {
	Consume(msgs <-chan amqp091.Delivery)
}

type PushNotificationConsumer struct {
	validator             *validator.Validate
	patientDataStore      datastore.PatientDataStore
	notificationDataStore datastore.NotificationDataStore
	notificationClient    notification.Client
	logger                *zap.SugaredLogger
}

func NewPushNotificationConsumer(validate *validator.Validate, patientDataStore datastore.PatientDataStore, notificationDataStore datastore.NotificationDataStore, notificationClient notification.Client, logger *zap.SugaredLogger) *PushNotificationConsumer {
	_ = validate.RegisterValidation("mapData", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(map[string]string)
		return ok
	})

	return &PushNotificationConsumer{
		validator:             validate,
		patientDataStore:      patientDataStore,
		notificationDataStore: notificationDataStore,
		notificationClient:    notificationClient,
		logger:                logger,
	}
}

func (c PushNotificationConsumer) Consume(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		go func() {
			isAck := c.Process(d)
			if isAck {
				c.assertError(d.Ack(false), "Failed to ack")
				return
			}
			c.assertError(d.Nack(false, true), "Failed to nack")
		}()
	}
}

func (c PushNotificationConsumer) Process(d amqp091.Delivery) bool {
	// Payload parsing and validation
	payload, err := c.ParsePayload(d.Body)
	if err != nil {
		c.logger.Warnw("message payload invalid form", "error", err.Error(), "payload", string(d.Body))
		return true
	}

	if err := c.validator.Struct(payload); err != nil {
		c.logger.Warnw("payload doesn't pass validation", "error", err.Error(), "payload", payload)
		return true
	}
	if err := c.validator.Var(payload.Data, "mapData"); err != nil {
		c.logger.Warnw("payload data doesn't pass validation", "error", err.Error(), "payload", payload)
		return true
	}

	// Query patient by given id in payload
	patient, err := c.patientDataStore.FindByIDOrRefID(payload.ID)
	if err != nil {
		c.assertError(err, "c.patientDataStore.FindByIDOrRefID error")
		return false
	}
	if patient == nil {
		c.logger.Warnw("patient not found", "id", payload.ID)
		return true
	}

	// Save notification in db
	noti := &datastore.Notification{
		Title:     payload.Title,
		Body:      payload.Body,
		IsRead:    false,
		PatientID: patient.ID,
	}
	err = c.notificationDataStore.Create(noti)
	if err != nil {
		c.assertError(err, "c.notificationDataStore.Create error")
		return false
	}

	if patient.NotificationToken == "" {
		return true
	}
	// Send push notification
	params := notification.SendParams{
		Token: patient.NotificationToken,
		Title: noti.Title,
		Body:  noti.Body,
	}
	payload.Data["notificationID"] = fmt.Sprintf("%d", noti.ID)
	err = c.notificationClient.Send(context.Background(), params, payload.Data)
	if err != nil {
		c.assertError(err, "c.notificationClient.Send error")
		return false
	}
	return true
}

func (c PushNotificationConsumer) ParsePayload(body []byte) (*Payload, error) {
	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (c PushNotificationConsumer) assertError(err error, msg string) {
	if err != nil {
		sentry.CaptureException(err)
		c.logger.Errorw(msg, "error", err.Error())
	}
}
