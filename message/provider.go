package message

import "net/url"

type Provider interface {
	Send(*Message) error
	Credit() (float64, error)
	Callback(url.Values) (string, error)
}
