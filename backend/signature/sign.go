package signature

import (
	"crypto"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func DefaultSignAlg() jwt.SigningMethod {
	return jwt.SigningMethodHS512
}

// SignedRequest creates JWT for a request
type SignedRequest struct {
	JWTPayload

	// Request payload (optional)
	Body io.Reader
	// Algorithm to use during body hashing
	BodyHashAlg crypto.Hash
}

func (r SignedRequest) Sign(smethod jwt.SigningMethod, privateKey interface{}, notBefore time.Time) (jwtToken string, err error) {
	switch {
	case smethod == nil:
		err = errors.New("signing method can't be empty")
		return
	case privateKey == nil:
		err = errors.New("private key can't be empty")
		return
	}

	// validate payload
	if err = newValidator().Struct(r.JWTPayload); err != nil {
		err = fmt.Errorf("invalid jwt payload: %w", err)
		return
	}

	// body (if presented) hash in hex form
	var (
		bodyHashStr string
		bodyHashAlg string
	)
	if r.Body != nil {
		if r.BodyHashAlg == 0 {
			err = errors.New("body reader is specified, but body hash alg is not")
			return
		}
		bodyHashStr, err = HashReader(r.Body, r.BodyHashAlg)
		if err != nil {
			err = fmt.Errorf("body hash: %w", err)
			return
		}
		bodyHashAlg = r.BodyHashAlg.String()
	}

	claims := JWTClaims{
		Method:      r.Method,
		URL:         r.URL,
		BodyHash:    bodyHashStr,
		BodyHashAlg: bodyHashAlg,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),

			Issuer:   r.Signer,
			Subject:  string(r.Subject),
			ID:       r.Nonce,
			Audience: r.Recipients,
		},
	}
	if !notBefore.IsZero() {
		claims.RegisteredClaims.NotBefore = jwt.NewNumericDate(notBefore.UTC())
	}

	// sign
	jwtToken, err = jwt.NewWithClaims(smethod, claims).SignedString(privateKey)
	if err != nil {
		err = fmt.Errorf("signing jwt: %w", err)
		return
	}
	return
}
