package notify

type ExpoPushMessage struct {
	To    string      `json:"to"`
	Sound string      `json:"sound"`
	Title string      `json:"title,omitempty"`
	Body  string      `json:"body"`
	Data  interface{} `json:"data,omitempty"`
}
