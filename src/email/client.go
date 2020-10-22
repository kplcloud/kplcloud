/**
 * @Time : 2019-07-18 15:05
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package email

import (
	"github.com/icowan/config"
	"net/http"
	"time"
)

type EmailInterface interface {
	AddAttachUrl(key, val string) EmailInterface
	AddInlineUrl(key, val string) EmailInterface
	SetTitle(title string) EmailInterface
	SetContent(content string) EmailInterface
	SetContentType(contentType string) EmailInterface
	AddEmailAddress(mail []string) EmailInterface
	AddCcEmailAddress(mail []string) EmailInterface
	AddBccEmailAddress(mail []string) EmailInterface
	SetHttpClient(client *http.Client) EmailInterface
	AddHeader(key, val string) EmailInterface
	AppointmentTime(t time.Time) EmailInterface
	Send() (err error)
}

func NewEmail(cf *config.Config) EmailInterface {

	// smtp 方式
	return NewEmailSmtp(&MailSmtpConf{
		User:     cf.GetString("email", "smtp_user"),
		Password: cf.GetString("email", "smtp_password"),
		Host:     cf.GetString("email", "smtp_host"),
	})
}
