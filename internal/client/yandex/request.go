package yandex

// Request модель запроса к АПИ ИИ
type Request struct {
	Model           string  `json:"model"`
	Instructions    string  `json:"instructions"`
	Input           string  `json:"input"`
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"max_output_tokens"`
}

func NewRequest(model, instructions, input string, temperature float64, maxTokens int) *Request {
	return &Request{
		Model:           model,
		Instructions:    instructions,
		Input:           input,
		Temperature:     temperature,
		MaxOutputTokens: maxTokens,
	}
}
