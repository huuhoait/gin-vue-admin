package response

type Email struct {
	To      string `json:"to"`      // email recipient
	Subject string `json:"subject"` // email subject
	Body    string `json:"body"`    // email body
}
