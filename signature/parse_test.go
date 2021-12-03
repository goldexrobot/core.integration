package signature

import (
	"bytes"
	"crypto"
	"errors"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestSignerRegexp(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		wantMatch bool
	}{
		{"short", "13", false},
		{"long", "123456789012345678901234567890123", false},
		{"wings", "-123456", false},
		{"wings", "123456-", false},
		{"ok", "123456", true},
		{"ok", "abc", true},
		{"ok", "1-Z", true},
		{"ok", "a-Z-0-9", true},
		{"specials", "login!", false},
		{"specials", "lo@in", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch := signerRex.MatchString(tt.str)
			if gotMatch != tt.wantMatch {
				t.Errorf("got %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestNonceRegexp(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		wantMatch bool
	}{
		{"short", "12345", false},
		{"ok", "123456", true},
		{"long", "c92143da-0c6e-4a27-bb65-1abda1958d8f0", false},
		{"ok", "12345678901234567890123456789012356", true},
		{"wings", "-c92143da-0c6e-4a27-bb65-1abda1958d8f", false},
		{"wings", "c92143da-0c6e-4a27-bb65-1abda1958d8f-", false},
		{"wings", "-c92143da-0c6e-4a27-bb65-1abda1958d8f-", false},
		{"ok", "c92143da-0c6e-4a27-bb65-1abda1958d8f", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch := nonceRex.MatchString(tt.str)
			if gotMatch != tt.wantMatch {
				t.Errorf("got %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestParseToken(t *testing.T) {

	var keypair = map[string]struct {
		Private []byte
		Public  []byte
	}{
		"goldex": {
			[]byte("awesomekey!"),
			[]byte("awesomekey!"),
		},
		"different-keys": {
			[]byte("foo"),
			[]byte("bar"),
		},
	}

	var sign = func(r SignedRequest, smethod jwt.SigningMethod) string {
		j, err := r.Sign(smethod, keypair[r.Signer].Private)
		if err != nil {
			t.Fatal(err)
		}
		return j
	}

	var l2k = func(signer string) (key interface{}, err error) {
		if pair, ok := keypair[signer]; ok {
			return pair.Public, nil
		}
		return nil, errors.New("no key for this login")
	}

	type args struct {
		jwtToken   string
		loginToKey func(login string) (key interface{}, err error)
	}
	tests := []struct {
		name    string
		args    args
		wantA   JWTPayload
		wantErr bool
	}{
		{
			"ok w/ body",
			args{
				sign(
					SignedRequest{
						JWTPayload: JWTPayload{
							Method:     "POST",
							URL:        "http://example.com",
							Nonce:      "aZ-1234567890",
							Signer:     "goldex",
							Subject:    "subject",
							Recipients: []string{"foo", "bar"},
						},
						Body:        bytes.NewBuffer([]byte("{}")),
						BodyHashAlg: crypto.SHA512,
					},
					DefaultSignAlg(),
				),
				l2k,
			},
			JWTPayload{
				Method:      "POST",
				URL:         "http://example.com",
				Nonce:       "aZ-1234567890",
				Signer:      "goldex",
				Subject:     "subject",
				Recipients:  []string{"foo", "bar"},
				BodyHash:    "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd",
				BodyHashAlg: "SHA-512",
			},
			false,
		},
		{
			"ok w/o body",
			args{
				sign(
					SignedRequest{
						JWTPayload: JWTPayload{
							Method: "GET",
							URL:    "http://example.com",
							Nonce:  "aZ-1234567890",
							Signer: "goldex",
						},
					},
					DefaultSignAlg(),
				),
				l2k,
			},
			JWTPayload{
				Method: "GET",
				URL:    "http://example.com",
				Nonce:  "aZ-1234567890",
				Signer: "goldex",
			},
			false,
		},
		{
			"invalid signature",
			args{
				sign(
					SignedRequest{
						JWTPayload: JWTPayload{
							Method: "GET",
							URL:    "http://example.com",
							Nonce:  "aZ-1234567890",
							Signer: "different-keys",
						},
					},
					DefaultSignAlg(),
				),
				l2k,
			},
			JWTPayload{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, err := ParseToken(tt.args.jwtToken, tt.args.loginToKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("ParseToken() = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}
