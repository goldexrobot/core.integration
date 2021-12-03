package signature

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
)

var (
	allowedBodyHashAlgs = map[string]crypto.Hash{
		crypto.SHA256.String():   crypto.SHA256,
		crypto.SHA384.String():   crypto.SHA384,
		crypto.SHA512.String():   crypto.SHA512,
		crypto.SHA3_224.String(): crypto.SHA3_224,
		crypto.SHA3_256.String(): crypto.SHA3_256,
		crypto.SHA3_384.String(): crypto.SHA3_384,
		crypto.SHA3_512.String(): crypto.SHA3_512,
	}
	allowedBodyHashAlgsNames []string
)

func init() {
	for k := range allowedBodyHashAlgs {
		allowedBodyHashAlgsNames = append(allowedBodyHashAlgsNames, k)
	}
}

func DefaultBodyHashAlg() crypto.Hash {
	return crypto.SHA512
}

// Hash returns a hash of a given bytes as solid hex string with no leading 0x.
// Example: "" => "", "{}" => "27c746...0763fd"
func Hash(body []byte, hasher crypto.Hash) (hx string, err error) {
	if len(body) == 0 {
		return
	}

	h := hasher.New()
	n, err := h.Write(body)
	if err != nil {
		return
	}
	if n != len(body) {
		err = fmt.Errorf("hash: wrote %v, expected %v", n, len(body))
		return
	}

	hx = hex.EncodeToString(h.Sum(nil))
	return
}

// HashReader is like Hash but hashes bytes from a reader
func HashReader(r io.Reader, hasher crypto.Hash) (hx string, err error) {
	if r == nil {
		return
	}

	h := hasher.New()
	n, err := io.Copy(h, r)
	if err != nil {
		return
	}
	if n == 0 {
		return
	}

	hx = hex.EncodeToString(h.Sum(nil))
	return
}
