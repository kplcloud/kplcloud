/**
 * @Time : 2019-07-18 15:05
 * @Author : solacowa@gmail.com
 * @File : smtp
 * @Software: GoLand
 */

package email

import (
	"bytes"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type MailSmtpConf struct {
	User, Password, Host string
}

type emailSmtp struct {
	user, password, host, subject, body, contentType string
	to,cc,bcc                                        []string

}

func NewEmailSmtp(conf *MailSmtpConf) EmailInterface {
	return &emailSmtp{
		user:     conf.User,
		password: conf.Password,
		host:     conf.Host,
	}
}

func (c *emailSmtp) AddAttachUrl(key, val string) EmailInterface {

	return c
}

func (c *emailSmtp) AddInlineUrl(key, val string) EmailInterface {

	return c
}

func (c *emailSmtp) SetTitle(title string) EmailInterface {
	c.subject = title
	return c
}

func (c *emailSmtp) SetContent(content string) EmailInterface {
	c.body = content
	return c
}

func (c *emailSmtp) SetContentType(contentType string) EmailInterface {
	c.contentType = contentType
	return c
}

func (c *emailSmtp) AddEmailAddress(mail []string) EmailInterface {
	c.to = mail
	return c
}

func (c *emailSmtp) AddCcEmailAddress(mail []string) EmailInterface {
	c.cc = mail
	return c
}

func (c *emailSmtp) AddBccEmailAddress(mail []string) EmailInterface {
	c.bcc = mail
	return c
}

func (c *emailSmtp) SetHttpClient(client *http.Client) EmailInterface {
	return c
}

func (c *emailSmtp) AddHeader(key, val string) EmailInterface {
	return c
}

func (c *emailSmtp) AppointmentTime(t time.Time) EmailInterface {
	return c
}

func (c *emailSmtp) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

func (c *emailSmtp) Send() (err error) {
	hp := strings.Split(c.host, ":")
	auth := smtp.PlainAuth("", c.user, c.password, hp[0])

	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary" //边界线
	Header := make(map[string]string)
	Header["From"] = c.user
	Header["To"] = strings.Join(c.to, ";")
	Header["Cc"] = strings.Join(c.cc, ";")
	Header["Bcc"] = strings.Join(c.bcc, ";")
	Header["Subject"] = c.subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	c.writeHeader(buffer, Header)

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + c.contentType + "\r\n"
	body += "\r\n" + c.body + "\r\n"
	buffer.WriteString(body)


	buffer.WriteString("\r\n--" + boundary + "--")

	to := append(c.to,c.cc...)
	to = append(to,c.bcc...)

	err = smtp.SendMail(c.host, auth, c.user, to, buffer.Bytes())
	return err
}



