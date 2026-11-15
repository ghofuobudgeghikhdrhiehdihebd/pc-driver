package schemas

type ChatHandlerMessage struct {
	Sender  string `json:"sender" msgpack:"sender"`
	Content string `json:"content" msgpack:"content"`
}

type ChatHandlerTool struct {
	Name        string            `json:"name" msgpack:"name"`
	Description string            `json:"description" msgpack:"description"`
	Parameters  map[string]string `json:"parameters" msgpack:"parameters"`
}

type ChatHandlerRequest struct {
	Model         string               `json:"model" msgpack:"model"`
	Messages      []ChatHandlerMessage `json:"messages" msgpack:"messages"`
	CurrentPrompt string               `json:"currentPrompt" msgpack:"currentPrompt"`
	Tools         []ChatHandlerTool    `json:"tools" msgpack:"tools"`
}
