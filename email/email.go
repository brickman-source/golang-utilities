/*
 * https://github.com/gcloudplatform/email/blob/master/email.go
 */

package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"mime"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strings"
)

type Mail struct {
	smtpServer string
	fromMail   *mail.Address
	password   string
}

type Attachment struct {
	Filename    string
	Data        []byte
	Inline      bool
	ContentType string
}

func New(smtpServer, password string, mailAddress *mail.Address) *Mail {
	return &Mail{
		smtpServer: smtpServer,
		fromMail:   mailAddress,
		password:   password,
	}
}

func (m *Mail) Send(title, body string, toEmail []*mail.Address) error {
	auth := smtp.PlainAuth(
		"",
		m.fromMail.Address,
		m.password,
		m.smtpServer,
	)

	to := ""
	var toEmails []string
	for _, e := range toEmail {
		to += "," + e.String()
		toEmails = append(toEmails, e.Address)
	}
	if to != "" {
		to = to[1:]
	}

	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.fromMail.String()))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", strings.Trim(mime.QEncoding.Encode("utf-8", title), "\"")))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")

	bodybase64 := base64.StdEncoding.EncodeToString([]byte(body))

	buf.WriteString(bodybase64)

	//fmt.Println(buf.String())

	err := smtp.SendMail(
		m.smtpServer,
		auth,
		m.fromMail.Address,
		toEmails,
		buf.Bytes(),
	)

	return err
}

// the email format see
// https://support.microsoft.com/zh-cn/kb/969854
// https://tools.ietf.org/html/rfc1341
func (m *Mail) SendWithAttachment(title, body string, toEmail []*mail.Address, attachment []*Attachment) error {
	auth := smtp.PlainAuth(
		"",
		m.fromMail.Address,
		m.password,
		m.smtpServer,
	)

	to := ""
	var toEmails []string
	for _, e := range toEmail {
		to += "," + e.String()
		toEmails = append(toEmails, e.Address)
	}
	if to != "" {
		to = to[1:]
	}

	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.fromMail.String()))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", strings.Trim(mime.QEncoding.Encode("utf-8", title), "\"")))
	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := ""

	boundary = genBoundary(28)
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary))
	buf.WriteString("This is a message with multiple parts in MIME format.\r\n")
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))

	bodybase64 := base64.StdEncoding.EncodeToString([]byte(body))

	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s\r\n", bodybase64))

	for _, attach := range attachment {
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		if attach.Inline {
			buf.WriteString("Content-Type: message/rfc822\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: inline; filename=\"%s\"\r\n\r\n", attach.Filename))
			buf.Write(attach.Data)
		} else {
			if attach.ContentType == "" {
				ext := filepath.Ext(attach.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					attach.ContentType = mimetype
				} else {
					attach.ContentType = "application/octet-stream"
				}
			}

			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", attach.ContentType))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attach.Filename))
			b := make([]byte, base64.StdEncoding.EncodedLen(len(attach.Data)))
			base64.StdEncoding.Encode(b, attach.Data)
			// write base64 content in lines of up to 76 chars
			for i, l := 0, len(b); i < l; i++ {
				buf.WriteByte(b[i])
				if (i+1)%76 == 0 {
					buf.WriteString("\r\n")
				}
			}
		}
	}
	buf.WriteString(fmt.Sprintf("\r\n--%s--", boundary))

	//fmt.Println(buf.String())

	err := smtp.SendMail(
		m.smtpServer,
		auth,
		m.fromMail.Address,
		toEmails,
		buf.Bytes(),
	)

	return err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const letterLen = 62

func genBoundary(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(letterLen)]
	}

	return string(b)
}
