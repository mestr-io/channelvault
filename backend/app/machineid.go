package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/denisbrodbeck/machineid"
)

func AppID(appKey string) string {
	id, _ := machineid.ID()
	return protect(appKey, id)
}

func protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
