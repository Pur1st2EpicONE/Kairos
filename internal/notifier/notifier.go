package notifier

import (
	"Kairos/internal/config"
	"Kairos/internal/models"
	"fmt"
	"net/http"
	"net/url"
)

type Notifier interface {
	Notify(notification models.Notification) error
}

func NewNotifier(config config.Notifier) Notifier {
	return newSender(config)
}

type Sender struct {
	telegramToken    string
	telegramReceiver string
}

func newSender(config config.Notifier) *Sender {
	return &Sender{
		telegramToken:    config.TelegramToken,
		telegramReceiver: config.TelegramReceiver,
	}
}

func (s *Sender) Notify(notification models.Notification) error {

	if err := s.sendTelegram(notification.Message); err != nil {
		return fmt.Errorf("unable to send Telegram notification: %w", err)
	}

	return nil

}

func (s *Sender) sendTelegram(message string) error {

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.telegramToken)

	data := url.Values{}
	data.Set("chat_id", s.telegramReceiver)
	data.Set("text", message)

	client := new(http.Client)

	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to POST form to Telegram API %s: %w", apiURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-OK status %s for chat_id %s", resp.Status, s.telegramReceiver)
	}

	return nil

}
