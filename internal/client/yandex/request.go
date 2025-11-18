package yandex

// Request модель запроса к АПИ ИИ
type Request struct {
	Model           string           `json:"model"`
	Messages        []RequestMessage `json:"messages"`
	Temperature     float64          `json:"temperature"`
	MaxOutputTokens int              `json:"max_output_tokens"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewRequest(model, instructions, input string, temperature float64, maxTokens int) *Request {
	return &Request{
		Model:       model,
		Temperature: temperature,
		Messages: []RequestMessage{
			{
				Role:    "system",
				Content: instructions,
			},
			{
				Role:    "user",
				Content: input,
			},
		},
		MaxOutputTokens: maxTokens,
	}
}
