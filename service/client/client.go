package client_service

import (
	"crypto/md5"
	"encoding/hex"
)

func CheckCertificate(clientBrand string) string {
	h := md5.New()
	h.Write([]byte(clientBrand))
	s := hex.EncodeToString(h.Sum(nil))
	return s
}
