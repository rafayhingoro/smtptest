package helper

import (
	"net/mail"

	"github.com/rafayhingoro/smtp2http/message"
)

func ExtractEmails(addr []*mail.Address, _ ...error) []string {
	ret := []string{}

	for _, e := range addr {
		ret = append(ret, e.Address)
	}

	return ret
}

func TransformStdAddressToEmailAddress(addr []*mail.Address) []*message.EmailAddress {
	ret := []*message.EmailAddress{}

	for _, e := range addr {
		ret = append(ret, &message.EmailAddress{
			Address: e.Address,
			Name:    e.Name,
		})
	}

	return ret
}

// func smtpsrvMesssage2EmailMessage(msg *smtpsrv.Context)
