package notification

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseNotificationClient struct {
	client *messaging.Client
}

type Config struct {
	FirebaseCredentialFilePath string `env:"FIREBASE_CRED_FILE_PATH,required"`
}

func NewFirebaseNotificationClient(ctx context.Context, cfg *Config) (*FirebaseNotificationClient, error) {
	opt := option.WithCredentialsFile(cfg.FirebaseCredentialFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &FirebaseNotificationClient{client: client}, nil
}

func (c FirebaseNotificationClient) Send(ctx context.Context, params SendParams, data map[string]string) error {
	_, err := c.client.Send(ctx, &messaging.Message{
		Token: params.Token,
		Notification: &messaging.Notification{
			Title: params.Title,
			Body:  params.Body,
		},
		Data: data,
	})
	return err
}
