package python_test

import (
	"encoding/base64"
	"encoding/hex"
)

func hexToBase64(s string) string {
	data, _ := hex.DecodeString(s)
	return base64.StdEncoding.EncodeToString(data)
}
