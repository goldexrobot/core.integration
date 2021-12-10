//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/goldexrobot/core.backend.integration/signature"
)

var stopper chan struct{}

func init() {
	stopper = make(chan struct{})
}

func main() {
	js.Global().Set("signature", js.ValueOf(
		map[string]interface{}{
			"sign":   js.FuncOf(Sign),
			"verify": js.FuncOf(Verify),
			"exit": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				close(stopper)
				return js.Undefined()
			}),
		},
	))

	fmt.Println("Signature library started")
	<-stopper
}

func Sign(this js.Value, args []js.Value) interface{} {
	defer recoverMe()

	if len(args) < 8 {
		panic("not enough args")
	}

	req := signature.SignedRequest{
		JWTPayload: signature.JWTPayload{
			Method:     args[0].String(),
			URL:        args[1].String(),
			Nonce:      args[2].String(),
			Signer:     args[3].String(),
			Subject:    args[4].String(),
			Recipients: []string{args[5].String()},
		},
	}

	b, err := getBytes(args[6])
	if err != nil {
		return result(nil, err)
	}
	if len(b) > 0 {
		req.Body = bytes.NewBuffer(b)
		req.BodyHashAlg = signature.DefaultBodyHashAlg()
	}

	token, err := req.Sign(signature.DefaultSignAlg(), []byte(args[7].String()))
	return result(token, err)
}

func Verify(this js.Value, args []js.Value) interface{} {
	defer recoverMe()

	if len(args) < 2 {
		panic("not enough args")
	}

	p, err := signature.ParseToken(args[0].String(), func(_ string) (publicKey interface{}, err error) {
		return []byte(args[1].String()), nil
	})
	if err != nil {
		return result(nil, err)
	}
	b, _ := json.MarshalIndent(p, "  ", "  ")
	return result(string(b), nil)
}

////// HELPERS //////

// Stops panic
func recoverMe() {
	if v := recover(); v != nil {
		s := ""
		switch x := v.(type) {
		case string:
			s = x
		case error:
			s = x.Error()
		default:
			s = "incorrect usage"
		}
		fmt.Println("Panic:", s)
	}
}

func result(i interface{}, err error) interface{} {
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	return map[string]interface{}{
		"result": i,
	}
}

// Copies Uint8Array to a new bytes buffer
func getBytes(v js.Value) ([]byte, error) {
	len := v.Length()
	if len == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, len)
	if n := js.CopyBytesToGo(buf, v); n != len {
		return nil, errors.New("failed to copy bytes")
	}
	return buf, nil
}

// Copies bytes to a new Uint8Array
func getUint8Array(b []byte) (js.Value, error) {
	if len(b) == 0 {
		return js.Null(), errors.New("empty buffer")
	}
	arr := js.Global().Get("Uint8Array").New(len(b))
	if n := js.CopyBytesToJS(arr, b); n != len(b) {
		return js.Null(), errors.New("failed to copy bytes")
	}
	return arr, nil
}
