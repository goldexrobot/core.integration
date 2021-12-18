all:
	GOOS=js GOARCH=wasm go build -ldflags='-w -s' -modfile ./../../../go.mod -o ./../../../docs/signature/signature.wasm main.go
	cp "$$GOROOT/misc/wasm/wasm_exec.js" ./../../../docs/signature/