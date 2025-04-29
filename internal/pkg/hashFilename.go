package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func HashFilename(filename string) string {
	hasher := sha256.New()
	hasher.Write([]byte(filename + time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))[:16]
}
