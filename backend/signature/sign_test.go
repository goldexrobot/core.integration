package signature

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestSignedRequest_Sign(t *testing.T) {
	key := []byte("deadbeef")

	type args struct {
		smethod    jwt.SigningMethod
		privateKey interface{}
	}

	argz := args{
		smethod:    DefaultSignAlg(),
		privateKey: key,
	}

	tests := []struct {
		name    string
		r       SignedRequest
		args    args
		wantErr bool
	}{
		{
			"ok",
			SignedRequest{
				JWTPayload: JWTPayload{
					Method:     "CONNECT",
					URL:        "http://example.com",
					Nonce:      "123456",
					Signer:     "goldex-robot",
					Subject:    "notification",
					Recipients: []string{"project"},
				},
			},
			argz,
			false,
		},
		{
			"ok",
			SignedRequest{
				JWTPayload: JWTPayload{
					Method:     "CONNECT",
					URL:        "http://example.com",
					Nonce:      "123456",
					Signer:     "goldex-robot",
					Subject:    "notification",
					Recipients: []string{"project"},
				},
				Body:        bytes.NewBuffer([]byte("{}")),
				BodyHashAlg: DefaultBodyHashAlg(),
			},
			argz,
			false,
		},
		{
			"no private key",
			SignedRequest{
				JWTPayload: JWTPayload{
					Method:     "CONNECT",
					URL:        "http://example.com",
					Nonce:      "123456",
					Signer:     "goldex-robot",
					Subject:    "notification",
					Recipients: []string{"project"},
				},
			},
			args{
				smethod:    DefaultSignAlg(),
				privateKey: nil,
			},
			true,
		},
		{
			"no sign alg",
			SignedRequest{
				JWTPayload: JWTPayload{
					Method:     "CONNECT",
					URL:        "http://example.com",
					Nonce:      "123456",
					Signer:     "goldex-robot",
					Subject:    "notification",
					Recipients: []string{"project"},
				},
			},
			args{
				smethod:    nil,
				privateKey: key,
			},
			true,
		},
		{
			"invalid body hash alg",
			SignedRequest{
				JWTPayload: JWTPayload{
					Method:     "CONNECT",
					URL:        "http://example.com",
					Nonce:      "123456",
					Signer:     "goldex-robot",
					Subject:    "notification",
					Recipients: []string{"project"},
				},
				Body:        bytes.NewBuffer([]byte("{}")),
				BodyHashAlg: 0,
			},
			argz,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Sign(tt.args.smethod, tt.args.privateKey, time.Now())
			if (err != nil) != tt.wantErr {
				t.Errorf("SignedRequest.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
