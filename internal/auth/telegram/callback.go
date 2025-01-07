package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/url"
	"slices"
	"strings"
)

// ValidateTelegramCallbackData https://core.telegram.org/widgets/login#checking-authorization
func (t *Telegram) ValidateTelegramCallbackData(queryMap url.Values) bool {
	dataParts := make([]string, 0, len(queryMap))

	var hash = ""
	for k, v := range queryMap {
		switch k {
		case "hash":
			hash = v[0]
		case "next":
			continue
		default:
			dataParts = append(dataParts, k+"="+v[0])
		}
	}

	slices.Sort(dataParts)

	imploded := strings.Join(dataParts, "\n")
	tokenHash := sha256.New()
	_, _ = io.WriteString(tokenHash, t.token)
	hmacTokenHash := hmac.New(sha256.New, tokenHash.Sum(nil))
	_, _ = io.WriteString(hmacTokenHash, imploded)
	ss := hex.EncodeToString(hmacTokenHash.Sum(nil))
	return hash == ss
}
