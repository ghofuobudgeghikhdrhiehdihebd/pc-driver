package schemas

type ResponseMessage struct {
	Content string `json:"content" msgpack:"content"`
}

type ToolCall struct {
	Name      string            `json:"name" msgpack:"name"`
	Arguments map[string]string `json:"arguments" msgpack:"arguments"`
}

type ChatHandlerResponse struct {
	Content   string           `json:"content,omitempty" msgpack:"content,omitempty"`
	Message   *ResponseMessage `json:"message,omitempty" msgpack:"message,omitempty"`
	ToolCalls []ToolCall       `json:"toolCalls,omitempty" msgpack:"toolCalls,omitempty"`
}
