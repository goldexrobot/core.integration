all:
	GOOS=js GOARCH=wasm go build -ldflags='-w -s' -o ./../../docs/signature/signature.wasm main.go
	cp "$$GOROOT/misc/wasm/wasm_exec.js" ./../../docs/signature/