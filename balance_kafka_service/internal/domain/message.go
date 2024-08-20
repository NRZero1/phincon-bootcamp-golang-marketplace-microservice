package domain

type Message[H any, B any] struct {
	Header H `json:"header"`
	Body   B `json:"body"`
}
