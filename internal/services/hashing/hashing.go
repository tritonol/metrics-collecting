package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func SignSHA256(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func VerifySHA256(data []byte, key, signature string) bool {
	expectedSignature := SignSHA256(data, key)
	return signature == expectedSignature
}
