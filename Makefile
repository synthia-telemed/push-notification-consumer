unit-test:
	ginkgo -r

mockgen:
	mockgen -source=pkg/notification/client.go -destination=test/mock_notification/mock_notification.go -package mock_notification
	mockgen -source=pkg/datastore/patient.go -destination=test/mock_datastore/mock_patient_datastore.go -package mock_datastore
	mockgen -source=pkg/datastore/notification.go -destination=test/mock_datastore/mock_notification.go -package mock_datastore