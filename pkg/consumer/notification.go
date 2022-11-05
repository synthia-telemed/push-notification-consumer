package consumer

import (
	"github.com/go-playground/validator/v10"
)

type PushNotificationMessageBody struct {
	Token string            `json:"token,omitempty" validate:"required"`
	Title string            `json:"title,omitempty" validate:"required"`
	Body  string            `json:"body,omitempty" validate:"required"`
	Data  map[string]string `json:"data,omitempty"`
}

func (b *PushNotificationMessageBody) IsValid() bool {
	validate := validator.New()
	err := validate.Struct(b)
	return err == nil
}
