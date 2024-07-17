package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"instagram-clone.com/m/pkg/logger"
)

const (
	CSRFHeader = "X-CSRF-Token"
	csrfSalt   = "HwI]7vfGR=>F/<,9JqeV0GS`m-j*A2jn6z4+c^4tBhQ#])af;V2ti-c@|r_7xkvb"
)

// Create CSRF token
func MakeToken(sid string, logger logger.Logger) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, csrfSalt+sid)
	if err != nil {
		logger.Errorf("Make CSRF Token", err)
	}
	token := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	return token
}

// Validate CSRF token
func ValidateToken(token string, sid string, logger logger.Logger) bool {
	trueToken := MakeToken(sid, logger)
	return token == trueToken
}
