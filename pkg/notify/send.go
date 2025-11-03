package notify

import (
	"bytes"
	"net/http"
	"subscriptionplus/server/config"
)

func SendPushNotification(token, title, message string) error {
	msg := ExpoPushMessage{
		To:    token,
		Sound: "default",
		Title: title,
		Body:  message,
	}
	body, _ := config.JSON.Marshal(msg)

	resp, err := http.Post("https://api.expo.dev/v2/push/send", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
