package propagators

import b64 "encoding/base64"

type Base64Crypto struct{}

func (_ Base64Crypto) Encrypt(data []byte) (out []byte, err error) {
	return []byte(b64.StdEncoding.EncodeToString(data)), nil
}

func (_ Base64Crypto) Decrypt(data []byte) (out []byte, err error) {
	return b64.StdEncoding.DecodeString(string(data))
}
