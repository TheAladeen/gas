package mail

import (
	"net/http"
	"net/smtp"

	"github.com/featen/ags/service/config"
	log "github.com/featen/ags/utils/log"
)

var auth smtp.Auth

func SendMail(receiver string, subject string, message string) int {
	if auth == nil {
		auth = smtp.PlainAuth(
			"",
			config.GetValue("SenderEmail"),
			config.GetValue("SenderPassword"),
			config.GetValue("SmtpServer"))
	}

	body := "To: " + receiver + "\r\nSubject: " + subject + "\r\n\r\n" + message
	err := smtp.SendMail(
		config.GetValue("SmtpServer")+":"+config.GetValue("SmtpPort"),
		auth,
		config.GetValue("SenderEmail"),
		[]string{receiver},
		[]byte(body))
	if err != nil {
		log.Info("mail %s sent to %s", subject, receiver)
		return http.StatusForbidden
	}
	return http.StatusOK
}
