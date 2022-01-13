package vars

import "flag"

var (
	FlagServerName     = flag.String("name", "mail.tld", "the server name")
	FlagListenAddr     = flag.String("listen", ":2525", "the smtp address to listen on")
	FlagWebhook        = flag.String("webhook", "http://localhost/webhook", "the webhook to send the data to")
	FlagMaxMessageSize = flag.Int64("msglimit", 1024*1024*2, "maximum incoming message size")
	FlagReadTimeout    = flag.Int("timeout.read", 5, "the read timeout in seconds")
	FlagWriteTimeout   = flag.Int("timeout.write", 5, "the write timeout in seconds")
)

func init() {
	flag.Parse()
}
