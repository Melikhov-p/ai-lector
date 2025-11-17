package yandex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/client"
)

// Client клиент взаимодействия с ИИ Яндекса
type Client struct {
	*client.Client
	OpenAIProject string
}

// NewClient новый клиент яндекса
func NewClient(name, token, address, model string) *Client {
	return &Client{
		Client:        client.NewClient(name, token, address, model),
		OpenAIProject: "b1gm2onh41hvp1g7idqg",
	}
}

// SetInstruction установить инструкции.
func (c *Client) SetInstruction(instruction string) {
	c.Instruction = instruction
}

// MakeRequest делает запрос к АПИ ИИ
func (c *Client) MakeRequest(req string) (string, error) {
	const op = "client.Yandex.MakeRequest"

	clientReq := NewRequest(c.Model, c.Instruction, req, 0.5, 500)

	body, _ := json.Marshal(clientReq)

	r, err := http.NewRequest("GET", req, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.Web.Do(r)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(respBody), nil
}
