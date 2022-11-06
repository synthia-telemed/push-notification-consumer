package consumer

type Payload struct {
	ID    string            `json:"id,omitempty" validate:"required"`
	Title string            `json:"title,omitempty" validate:"required"`
	Body  string            `json:"body,omitempty" validate:"required"`
	Data  map[string]string `json:"data,omitempty"`
}
