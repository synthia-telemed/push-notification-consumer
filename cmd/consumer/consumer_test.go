package consumer_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rabbitmq/amqp091-go"
	Consumer "github.com/synthia-telemed/push-notification-consumer/cmd/consumer"
	"github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
	Notification "github.com/synthia-telemed/push-notification-consumer/pkg/notification"
	testhelper "github.com/synthia-telemed/push-notification-consumer/test/helper"
	"github.com/synthia-telemed/push-notification-consumer/test/mock_datastore"
	"github.com/synthia-telemed/push-notification-consumer/test/mock_notification"
	"go.uber.org/zap"
	"math/rand"
)

var _ = Describe("Push notification consumer", func() {
	var (
		consumer *Consumer.PushNotificationConsumer

		mockCtrl                  *gomock.Controller
		mockPatientDataStore      *mock_datastore.MockPatientDataStore
		mockNotificationDataStore *mock_datastore.MockNotificationDataStore
		mockNotificationClient    *mock_notification.MockClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockNotificationDataStore = mock_datastore.NewMockNotificationDataStore(mockCtrl)
		mockNotificationClient = mock_notification.NewMockClient(mockCtrl)
		consumer = Consumer.NewPushNotificationConsumer(validator.New(), mockPatientDataStore, mockNotificationDataStore, mockNotificationClient, zap.NewNop().Sugar())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Process", func() {
		var (
			payload Consumer.Payload
			d       amqp091.Delivery
			isAck   bool

			patient      *datastore.Patient
			notification *datastore.Notification
			sendParams   Notification.SendParams
		)

		BeforeEach(func() {
			payload = Consumer.Payload{
				ID:    uuid.NewString(),
				Title: uuid.NewString(),
				Body:  uuid.NewString(),
				Data:  map[string]string{"test": uuid.NewString()},
			}
			body, err := json.Marshal(&payload)
			Expect(err).To(BeNil())
			d = amqp091.Delivery{Body: body}

			patient = &datastore.Patient{RefID: uuid.NewString(), ID: uint(rand.Uint32()), NotificationToken: uuid.NewString()}
			notification = &datastore.Notification{Title: payload.Title, Body: payload.Body, IsRead: false, PatientID: patient.ID}
			payload.Data["notificationID"] = fmt.Sprintf("%d", notification.ID)
			sendParams = Notification.SendParams{Title: payload.Title, Body: payload.Body, Token: patient.NotificationToken}
		})

		JustBeforeEach(func() {
			isAck = consumer.Process(d)
		})

		When("payload is no valid JSON", func() {
			BeforeEach(func() {
				d = amqp091.Delivery{
					Body: []byte("awd"),
				}
			})
			It("should return isAck as true", func() {
				Expect(isAck).To(BeTrue())
			})
		})
		When("payload is invalid", func() {
			BeforeEach(func() {
				d = amqp091.Delivery{
					Body: []byte(`{"title": "test", "body": 123}`),
				}
			})
			It("should return isAck as true", func() {
				Expect(isAck).To(BeTrue())
			})
		})
		When("payload data is invalid", func() {
			BeforeEach(func() {
				d = amqp091.Delivery{
					Body: []byte(`{"id": "1", "title": "test", "body": "123", "data": "not-map"}`),
				}
			})
			It("should return isAck as true", func() {
				Expect(isAck).To(BeTrue())
			})
		})

		When("find patient by id or ref_id error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByIDOrRefID(payload.ID).Return(nil, testhelper.MockError)
			})
			It("should return isAck as false", func() {
				Expect(isAck).To(BeFalse())
			})
		})
		When("patient is not found", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByIDOrRefID(payload.ID).Return(nil, nil)
			})
			It("should return isAck as true", func() {
				Expect(isAck).To(BeTrue())
			})
		})

		When("save notification to db error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByIDOrRefID(payload.ID).Return(patient, nil)
				mockNotificationDataStore.EXPECT().Create(notification).Return(testhelper.MockError)
			})
			It("should return isAck as false", func() {
				Expect(isAck).To(BeFalse())
			})
		})

		When("patient's notification token is empty", func() {
			BeforeEach(func() {
				patient.NotificationToken = ""
				mockNotificationDataStore.EXPECT().Create(notification)
				mockPatientDataStore.EXPECT().FindByIDOrRefID(payload.ID).Return(patient, nil)
			})
			It("should return isAck as true", func() {
				Expect(isAck).To(BeTrue())
			})
		})
		When("send push notification error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByIDOrRefID(payload.ID).Return(patient, nil)
				mockNotificationDataStore.EXPECT().Create(notification).Return(nil)
				mockNotificationClient.EXPECT().Send(gomock.Any(), sendParams, payload.Data).Return(testhelper.MockError)
			})
			It("should return isAck as false", func() {
				Expect(isAck).To(BeFalse())
			})
		})

		When("no error occurred", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByIDOrRefID(payload.ID).Return(patient, nil)
				mockNotificationDataStore.EXPECT().Create(notification).Return(nil)
				mockNotificationClient.EXPECT().Send(gomock.Any(), sendParams, payload.Data).Return(nil)
			})
			It("should return isAck as true", func() {
				Expect(isAck).To(BeTrue())
			})
		})
	})
})
