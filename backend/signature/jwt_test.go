package signature

import (
	"crypto"
	"testing"
)

func TestJWTPayloadValidation(t *testing.T) {
	var validator = newValidator()

	var def = func(modify func(p *JWTPayload)) JWTPayload {
		p := JWTPayload{
			Method:      "GET",
			URL:         "http://example.com",
			Nonce:       "1234567890",
			Signer:      "signer1",
			Subject:     "subject666",
			Recipients:  []string{"recipient1", "recipient2"},
			BodyHash:    "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd",
			BodyHashAlg: crypto.SHA512.String(),
		}
		modify(&p)
		return p
	}

	tests := []struct {
		name      string
		payload   JWTPayload
		wantError bool
	}{
		{
			"ok",
			def(func(p *JWTPayload) {}),
			false,
		},
		{
			"ok / no optionals",
			def(func(p *JWTPayload) {
				p.Subject = ""
				p.Recipients = nil
				p.BodyHash = ""
				p.BodyHashAlg = ""
			}),
			false,
		},
		{
			"ok / no body",
			def(func(p *JWTPayload) {
				p.BodyHash = ""
				p.BodyHashAlg = ""
			}),
			false,
		},
		{
			"no method",
			def(func(p *JWTPayload) {
				p.Method = ""
			}),
			true,
		},
		{
			"invalid nonce",
			def(func(p *JWTPayload) {
				p.Nonce = "-a-zA-Z0-9"
			}),
			true,
		},
		{
			"short nonce",
			def(func(p *JWTPayload) {
				p.Nonce = "abc"
			}),
			true,
		},
		{
			"long nonce",
			def(func(p *JWTPayload) {
				p.Nonce = "1234567890123456789012345678901234567"
			}),
			true,
		},
		{
			"long method",
			def(func(p *JWTPayload) {
				p.Method = "123456789"
			}),
			true,
		},
		{
			"invalid nonce",
			def(func(p *JWTPayload) {
				p.Nonce = "12345!67890"
			}),
			true,
		},
		{
			"invalid signer",
			def(func(p *JWTPayload) {
				p.Signer = ""
			}),
			true,
		},
		{
			"invalid url",
			def(func(p *JWTPayload) {
				p.URL = "localhost"
			}),
			true,
		},
		{
			"invalid subject",
			def(func(p *JWTPayload) {
				p.Subject = "special!"
			}),
			true,
		},
		{
			"long subject",
			def(func(p *JWTPayload) {
				p.Subject = "123456789012345678901234567890123"
			}),
			true,
		},
		{
			"short body hash",
			def(func(p *JWTPayload) {
				p.BodyHash = "27c74670"
			}),
			true,
		},
		{
			"long body hash",
			def(func(p *JWTPayload) {
				p.BodyHash = "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd0000000000"
			}),
			true,
		},
		{
			"non-hex body hash",
			def(func(p *JWTPayload) {
				p.BodyHash = "z27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9"
			}),
			true,
		},
		{
			"no body alg with body hash",
			def(func(p *JWTPayload) {
				p.BodyHashAlg = ""
			}),
			true,
		},
		{
			"invalid body hash alg",
			def(func(p *JWTPayload) {
				p.BodyHashAlg = "FOO-512"
			}),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Struct(tt.payload)
			if (got != nil) != tt.wantError {
				t.Errorf("Struct() = %v, want %v", got, tt.wantError)
			}
		})
	}
}
