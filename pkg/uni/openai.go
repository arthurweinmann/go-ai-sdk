package uni

type OpenAIWrapper struct {
	Temperature float32
}

func (w *OpenAIWrapper) SetTemperature(t float32) *OpenAIWrapper {
	w.Temperature = t
	return w
}

func (w *OpenAIWrapper) AddMessage(ChatMessage) error {
	return nil
}
