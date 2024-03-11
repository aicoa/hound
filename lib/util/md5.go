package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5Value(context []byte) string {
	hash := md5.New()
	io.WriteString(hash, string(context))
	return hex.EncodeToString(hash.Sum(nil))
}
