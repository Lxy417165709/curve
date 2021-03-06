package utils

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"html/template"
)

type VrcEmailSender struct {
	emailClient   *EmailClient
	vrcGenerator  *RandStringGenerator
	emailTemplate *template.Template
}

func NewVrcEmailSender(client *EmailClient, generator *RandStringGenerator, emailTemplate *template.Template) *VrcEmailSender {
	return &VrcEmailSender{
		emailClient:   client,
		vrcGenerator:  generator,
		emailTemplate: emailTemplate,
	}
}

func (v *VrcEmailSender) SendVrcEmail(subject string, emailAddrOfReceiver string, vrcExpiredSecond int) (string, error) {
	var htmlContentBuffer bytes.Buffer
	vrc := v.vrcGenerator.Get()
	if err := v.emailTemplate.Execute(&htmlContentBuffer, struct {
		Vrc            string
		VrcExpiredSecond int
	}{
		Vrc:            vrc,
		VrcExpiredSecond:vrcExpiredSecond,
	}); err != nil {
		logs.Error(err)
		return "", err
	}
	if err := v.emailClient.SendEmail(emailAddrOfReceiver, subject, htmlContentBuffer.String()); err != nil {
		logs.Error(err)
		return "", err
	}
	return vrc, nil
}
