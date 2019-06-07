package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex" 
	"fmt"
)

func TokenGenerate() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func GetMD5(text string) string {
	h := md5.New()
    h.Write([]byte(text))
    cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
