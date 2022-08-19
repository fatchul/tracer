package attribute

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

type (
	HeaderAttr struct {
		Referrer  string `json:"header.referer"`
		Origin    string `json:"header.origin"`
		CountryID string `json:"header.country_id"`
		Timezone  string `json:"header.timezone"`
		ClientID  string `json:"header.client_id"`
		Platform  string `json:"header.platform"`
	}
)

func Header(req *http.Request) []attribute.KeyValue {
	return parseAttrToKV(HeaderAttr{
		Referrer:  req.Header.Get("referer"),
		Origin:    req.Header.Get("origin"),
		CountryID: req.Header.Get("country-id"),
		Timezone:  req.Header.Get("timezone"),
		ClientID:  req.Header.Get("client-id"),
		Platform:  req.Header.Get("sec-ch-ua-platform"),
	})
}
