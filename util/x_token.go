package util

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const AppPrefix = "SUBPLUS"

func GenerateXToken(userUUID string, planID int) string {
	rand.Seed(time.Now().UnixNano())

	salt := fmt.Sprintf("%04d", rand.Intn(9999))
	base := fmt.Sprintf("%s-%d-%s", userUUID, planID, salt)

	hash := sha256.Sum256([]byte(base))

	tokenPart := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:])

	tokenPart = strings.ToLower(tokenPart[:10])

	token := fmt.Sprintf("%s_%s", AppPrefix, tokenPart)

	return token
}
