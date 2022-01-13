package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"time"

	"github.com/alash3al/go-smtpsrv/v3"
	"github.com/go-resty/resty/v2"
	"github.com/rafayhingoro/smtp2http/helper"
	"github.com/rafayhingoro/smtp2http/message"
	"github.com/rafayhingoro/smtp2http/vars"
)

func main() {
	cfg := smtpsrv.ServerConfig{
		ReadTimeout:     time.Duration(*vars.FlagReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(*vars.FlagWriteTimeout) * time.Second,
		ListenAddr:      *vars.FlagListenAddr,
		MaxMessageBytes: int(*vars.FlagMaxMessageSize),
		BannerDomain:    *vars.FlagServerName,
		Handler: smtpsrv.HandlerFunc(func(c *smtpsrv.Context) error {
			msg, err := c.Parse()
			if err != nil {
				return errors.New("Cannot read your message: " + err.Error())
			}

			spfResult, _, _ := c.SPF()

			jsonData := message.EmailMessage{
				ID:            msg.MessageID,
				Date:          msg.Date.String(),
				References:    msg.References,
				SPFResult:     spfResult.String(),
				ResentDate:    msg.ResentDate.String(),
				ResentID:      msg.ResentMessageID,
				Subject:       msg.Subject,
				Attachments:   []*message.EmailAttachment{},
				EmbeddedFiles: []*message.EmailEmbeddedFile{},
			}

			jsonData.Body.HTML = string(msg.HTMLBody)
			jsonData.Body.Text = string(msg.TextBody)

			jsonData.Addresses.From = helper.TransformStdAddressToEmailAddress([]*mail.Address{c.From()})[0]
			jsonData.Addresses.To = helper.TransformStdAddressToEmailAddress([]*mail.Address{c.To()})[0]

			jsonData.Addresses.Cc = helper.TransformStdAddressToEmailAddress(msg.Cc)
			jsonData.Addresses.Bcc = helper.TransformStdAddressToEmailAddress(msg.Bcc)
			jsonData.Addresses.ReplyTo = helper.TransformStdAddressToEmailAddress(msg.ReplyTo)
			jsonData.Addresses.InReplyTo = msg.InReplyTo

			if resentFrom := helper.TransformStdAddressToEmailAddress(msg.ResentFrom); len(resentFrom) > 0 {
				jsonData.Addresses.ResentFrom = resentFrom[0]
			}

			jsonData.Addresses.ResentTo = helper.TransformStdAddressToEmailAddress(msg.ResentTo)
			jsonData.Addresses.ResentCc = helper.TransformStdAddressToEmailAddress(msg.ResentCc)
			jsonData.Addresses.ResentBcc = helper.TransformStdAddressToEmailAddress(msg.ResentBcc)

			for _, a := range msg.Attachments {
				data, _ := ioutil.ReadAll(a.Data)
				jsonData.Attachments = append(jsonData.Attachments, &message.EmailAttachment{
					Filename:    a.Filename,
					ContentType: a.ContentType,
					Data:        base64.StdEncoding.EncodeToString(data),
				})
			}

			for _, a := range msg.EmbeddedFiles {
				data, _ := ioutil.ReadAll(a.Data)
				jsonData.EmbeddedFiles = append(jsonData.EmbeddedFiles, &message.EmailEmbeddedFile{
					CID:         a.CID,
					ContentType: a.ContentType,
					Data:        base64.StdEncoding.EncodeToString(data),
				})
			}

			resp, err := resty.New().R().SetHeader("Content-Type", "application/json").SetBody(jsonData).Post(*vars.FlagWebhook)
			if err != nil {
				log.Println(err)
				return errors.New("E1: Cannot accept your message due to internal error, please report that to our engineers")
			} else if resp.StatusCode() != 200 {
				log.Println(resp.Status())
				return errors.New("E2: Cannot accept your message due to internal error, please report that to our engineers")
			}

			return nil
		}),
	}

	fmt.Println(smtpsrv.ListenAndServe(&cfg))
}
