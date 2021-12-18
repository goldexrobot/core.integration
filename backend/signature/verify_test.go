package signature

import (
	"crypto"
	"net/url"
	"testing"
)

func TestEqualURL(t *testing.T) {
	must := func(s string) *url.URL {
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		return u
	}

	type args struct {
		a *url.URL
		b *url.URL
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"same", args{must("http://a.com"), must("http://a.com")}, true},
		{"same / ignore scheme", args{must("http://a.com"), must("https://a.com")}, true},
		{"same / case", args{must("http://a.com"), must("http://A.com")}, true},
		{"same / fragment", args{must("http://a.com#x"), must("http://a.com#y")}, true},
		{"diff", args{must("http://a.com"), must("http://b.com")}, false},

		{"same / port", args{must("http://a.com:8080"), must("http://a.com:8080")}, true},
		{"diff / port", args{must("http://a.com:8080"), must("http://a.com:8081")}, false},

		{"same / path", args{must("http://a.com/path/to"), must("http://a.com/path/to")}, true},
		{"same / path case", args{must("http://a.com/path/to"), must("http://a.com/path/TO")}, true},
		{"same / redundant path", args{must("http://a.com"), must("http://a.com/")}, true},
		{"same / redundant nonempty path", args{must("http://a.com/path/to"), must("http://a.com/path/to/")}, true},
		{"diff / path", args{must("http://a.com/path/to"), must("http://a.com/path/to/something")}, false},

		{"same / redundant query", args{must("http://a.com"), must("http://a.com?")}, true},
		{"same / query", args{must("http://a.com?a=b&c=d"), must("http://a.com?a=b&c=d")}, true},
		{"same / query unsorted", args{must("http://a.com?a=b&c=d"), must("http://a.com?c=d&a=b")}, true},
		{"same / query empty key", args{must("http://a.com?a"), must("http://a.com?a=")}, true},
		{"diff / query", args{must("http://a.com?a=b"), must("http://a.com?a")}, false},
		{"diff / query key case", args{must("http://a.com?a=b"), must("http://a.com?A=b")}, false},
		{"diff / query value case", args{must("http://a.com?a=b"), must("http://a.com?a=B")}, false},

		{"same / query + as space", args{must("http://a.com?a=f%20oo"), must("http://a.com?a=f+oo")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualURL(tt.args.a, tt.args.b, false); got != tt.want {
				t.Errorf("EqualURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyRequest(t *testing.T) {
	must := func(s string) *url.URL {
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		return u
	}

	type args struct {
		method string
		urll   *url.URL
		t      JWTPayload
		body   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ok w/o body",
			args{
				method: "GET",
				urll:   must("http://example.com"),
				t: JWTPayload{
					Method: "get",
					URL:    "http://example.com/?#",
				},
			},
			false,
		},
		{
			"ok w/ body",
			args{
				method: "POST",
				urll:   must("http://example.com/path/to/something?foo=bar&baz=qux"),
				t: JWTPayload{
					Method:      "PoSt",
					URL:         "http://example.com/path/to/something/?baz=qux&foo=bar#",
					BodyHash:    "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd",
					BodyHashAlg: crypto.SHA512.String(),
				},
				body: []byte("{}"),
			},
			false,
		},
		{
			"hash for empty body",
			args{
				method: "POST",
				urll:   must("http://example.com/path/to/something?foo=bar&baz=qux"),
				t: JWTPayload{
					Method:      "PoSt",
					URL:         "http://example.com/path/to/something/?baz=qux&foo=bar#",
					BodyHash:    "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd",
					BodyHashAlg: crypto.SHA512.String(),
				},
			},
			true,
		},
		{
			"no hash for body",
			args{
				method: "POST",
				urll:   must("http://example.com/path/to/something?foo=bar&baz=qux"),
				t: JWTPayload{
					Method: "PoSt",
					URL:    "http://example.com/path/to/something/?baz=qux&foo=bar#",
				},
				body: []byte("{}"),
			},
			true,
		},
		{
			"invalid hash",
			args{
				method: "POST",
				urll:   must("http://example.com/path/to/something?foo=bar&baz=qux"),
				t: JWTPayload{
					Method:      "PoSt",
					URL:         "http://example.com/path/to/something/?baz=qux&foo=bar#",
					BodyHash:    "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd",
					BodyHashAlg: crypto.SHA512.String(),
				},
				body: []byte("[]"),
			},
			true,
		},
		{
			"method mismatched",
			args{
				method: "GET",
				urll:   must("http://example.com"),
				t: JWTPayload{
					Method: "POST",
					URL:    "http://example.com/?#",
				},
			},
			true,
		},
		{
			"url mismatched",
			args{
				method: "GET",
				urll:   must("http://example.com?foo=bar"),
				t: JWTPayload{
					Method: "POST",
					URL:    "http://example.com?foo=baz",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyRequest(tt.args.method, tt.args.urll, tt.args.t, true, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("VerifyRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
