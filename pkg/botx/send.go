package botx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func Send(body string) {
	data := map[string]string{
		"chat_id":    os.Getenv("BOT_CHAT_ID"),
		"text":       body,
		"parse_mode": "MarkdownV2",
	}

	jsonBody, _ := json.Marshal(data)

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("BOT_TOKEN")),
		"application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("[ERROR] Failed to send Telegram message:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("[INFO] Send bot new transaction, status:", resp.Status)
}
