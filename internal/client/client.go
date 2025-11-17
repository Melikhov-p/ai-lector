package client

import (
	"net/http"
	"time"
)

// Client базовый класс ИИ клиента который ходит за запросами
type Client struct {
	Name        string
	Token       string
	Address     string // https://llm.api.cloud.yandex.net/v1/chat/completions
	Model       string // gpt://b1gm2onh41hvp1g7idqg/yandexgpt/latest
	Instruction string
	Web         *http.Client
}

// NewClient новый клиент
func NewClient(name, token, address, model string) *Client {
	wc := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       15 * time.Second,
	}

	return &Client{
		Name:    name,
		Token:   token,
		Address: address,
		Model:   model,
		Web:     &wc,
	}
}
