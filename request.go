package main

type PublishMessageRequest struct {
	Subject string      `json:"subject" validate:"required"`
	Payload interface{} `json:"payload" validate:"required"`
}
