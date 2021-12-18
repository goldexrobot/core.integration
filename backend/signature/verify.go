package signature

import (
	"fmt"
	"net/url"
	"strings"
)

// VerifyRequest verifies parsed and validated token vs actual HTTP request
func VerifyRequest(method string, urll *url.URL, t JWTPayload, compareHosts bool, body []byte) (err error) {
	tu, err := url.Parse(t.URL)
	if err != nil {
		err = fmt.Errorf("failed to parse url from token")
		return
	}

	// same method
	if !strings.EqualFold(t.Method, method) {
		err = fmt.Errorf("http method mismatched: signed %q, actual %q", t.Method, method)
		return
	}

	// almost same url
	if !EqualURL(tu, urll, !compareHosts) {
		err = fmt.Errorf("url mismatched: signed %q, actual %q", tu.String(), urll.String())
		return
	}

	// hash body
	var bodyHashStr string
	if len(body) > 0 {
		alg, ok := allowedBodyHashAlgs[t.BodyHashAlg]
		if !ok {
			err = fmt.Errorf("body hash alg %q not implemented", t.BodyHashAlg)
			return
		}

		bodyHashStr, err = Hash(body, alg)
		if err != nil {
			err = fmt.Errorf("body hash: %w", err)
			return
		}
	}

	// compare body hash
	if t.BodyHash != bodyHashStr {
		err = fmt.Errorf("body hash mismatched: signed %q, actual %q", t.BodyHash, bodyHashStr)
		return
	}

	return
}

// EqualURL checks url equality respecting only: user/password, host, port, path, query
func EqualURL(a, b *url.URL, skipHost bool) bool {
	return a.User.String() == b.User.String() &&
		(skipHost || strings.EqualFold(a.Host, b.Host)) &&
		strings.EqualFold(strings.Trim(a.Path, "/"), strings.Trim(b.Path, "/")) &&
		a.Query().Encode() == b.Query().Encode()
}
