package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"website/backend/telegram/utils"
)

func SendMessageGmail(name, email, message string) error {
	utils.LoadEnv()

	ctx := context.Background()
	gmailService, err := gmail.NewService(ctx)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %v", err)
	}

	msg := createMessage(name, email, message)

	_, err = gmailService.Users.Messages.Send("me", &gmail.Message{
		Raw: msg,
	}).Do()
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %v", err)
	}

	return nil
}

func createMessage(name, email, message string) string {
	body := fmt.Sprintf("From: %s\r\n", email)
	body += fmt.Sprintf("To: %s\r\n", "user@example.com") // Замените на адрес получателя
	body += fmt.Sprintf("Subject: %s\r\n", "Имя:"+name)
	body += "\r\n" + message

	raw := base64.URLEncoding.EncodeToString([]byte(body))

	return raw
}
