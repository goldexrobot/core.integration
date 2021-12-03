package signature

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

// JWTPayload is abstraction of what JWT transfers
type JWTPayload struct {
	// HTTP method
	Method string `validate:"required,min=1,max=8,alpha"`
	// Full URL
	URL string `validate:"required,url,max=256"`
	// Unique request ID/nonce
	Nonce string `validate:"required,nonce"`
	// Signer identity
	Signer string `validate:"required,signer"`
	// Subject of the request (optional)
	Subject string `validate:"omitempty,max=32,alphanum"`
	// Destination recipient (optional)
	Recipients []string `validate:"dive,required,signer"`

	// Body hash and algorithm (empty for bodyless requests).
	BodyHash    string `validate:"required_with=BodyHashAlg,omitempty,min=32,max=128,hexadecimal"`
	BodyHashAlg string `validate:"required_with=BodyHash,omitempty,min=1,max=16,bodyhashalg"`
}

// Extended JWT claims
type JWTClaims struct {
	BodyHashAlg string `json:"bha,omitempty"`
	BodyHash    string `json:"bhs,omitempty"`
	Method      string `json:"mtd,omitempty"`
	URL         string `json:"url,omitempty"`
	jwt.RegisteredClaims
}

////// VALIDATION //////

var (
	signerRex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]{1,30}[a-zA-Z0-9]$`) // 3 to 32 symbols
	nonceRex  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{4,34}[a-zA-Z0-9]$`)  // 6 to 36 symbols
)

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("signer", validateSigner)
	v.RegisterValidation("nonce", validateNonce)
	v.RegisterValidation("bodyhashalg", validateBodyHashAlg)
	return v
}

func validateSigner(fl validator.FieldLevel) bool {
	return signerRex.MatchString(fl.Field().String())
}
func validateNonce(fl validator.FieldLevel) bool {
	return nonceRex.MatchString(fl.Field().String())
}
func validateBodyHashAlg(fl validator.FieldLevel) bool {
	alg := strings.ToUpper(fl.Field().String())
	for _, a := range allowedBodyHashAlgsNames {
		if alg == a {
			return true
		}
	}
	return false
}
