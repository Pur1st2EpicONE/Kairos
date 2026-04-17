// Package notifier provides a simple notification system that can send messages
// via different channels. Currently, only Telegram is implemented.
package notifier

import (
	"Kairos/internal/config"
	"Kairos/internal/models"
	"fmt"
	"net/http"
	"net/url"
)

// Notifier defines the interface for sending notifications.
// Implementations should handle the delivery of a Notification to the appropriate channel.
type Notifier interface {
	Notify(notification models.Notification) error // Notify sends the given notification through the channel specified in the struct.
}

// NewNotifier constructs a Notifier instance based on the provided configuration.
// Currently, it returns a Telegram sender.
func NewNotifier(config config.Notifier) Notifier {
	return newSender(config)
}

// Sender is the concrete implementation of Notifier using Telegram.
// It holds the bot token and the receiver's chat ID.
type Sender struct {
	telegramToken    string // Bot token for Telegram API
	telegramReceiver string // Chat ID or username that receives messages
}

// newSender creates a new Sender with the given configuration.
func newSender(config config.Notifier) *Sender {
	return &Sender{
		telegramToken:    config.TelegramToken,
		telegramReceiver: config.TelegramReceiver,
	}
}

// Notify sends a notification via Telegram.
// It extracts the message from the notification and calls sendTelegram.
// If the Telegram send fails, it returns a wrapped error.
func (s *Sender) Notify(notification models.Notification) error {

	if err := s.sendTelegram(notification.Message); err != nil {
		return fmt.Errorf("unable to send Telegram notification: %w", err)
	}

	return nil

}

// sendTelegram sends a plain text message to the configured Telegram chat ID
// using the bot token. It constructs a POST request to the Telegram Bot API
// and checks the response status. Returns an error if the request fails or
// the API returns a non-OK status.
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
