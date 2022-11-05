package consumer

import (
	"github.com/go-playground/validator/v10"
)

type PushNotificationMessageBody struct {
	Token string            `json:"token" validate:"required"`
	Title string            `json:"title" validate:"required"`
	Body  string            `json:"body" validate:"required"`
	Data  map[string]string `json:"data"`
}

func (b *PushNotificationMessageBody) IsValid() bool {
	validate := validator.New()
	err := validate.Struct(b)
	_, isError := err.(validator.ValidationErrors)
	return isError
}
