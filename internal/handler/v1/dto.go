package v1

type RegisterDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateEventDTO struct {
	ID          string `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	TotalSeats  int    `json:"total_seats"`
}

type CreateNotificationV1 struct {
	Channel string   `json:"channel"` // The channel to send the notification through (e.g., "email", "telegram").
	Subject string   `json:"subject"` // The subject or title of the notification (used for email, optional for other channels).
	Message string   `json:"message"` // The main content of the notification.
	SendAt  string   `json:"send_at"` // The scheduled send time in RFC3339 format.
	SendTo  []string `json:"send_to"` // The list of recipients for the notification.
}
