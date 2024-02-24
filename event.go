package main

type Payload struct {
	ChatMessage string            `json:"chat_message"`
	Username    string            `json:"username"`
	Headers     map[string]string `json:"HEADERS"`
	TimeSent    string            `json:"time"`
}
