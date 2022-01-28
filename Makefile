all: backend terminal
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: backend terminal

backend:
	$(info Backend:)
	@$(MAKE) --no-print-directory -C backend

terminal:
	$(info Terminal:)
	@$(MAKE) --no-print-directory -C terminal

emu:
	@mkdir -p ./dist
	GOOS=windows GOARCH=amd64 go build -tags "-w -s" -o ./dist/tapi-emu-windows.exe ./terminal/cmd/api-emulator/main.go
	GOOS=linux GOARCH=amd64 go build -tags "-w -s" -o ./dist/tapi-emu-linux ./terminal/cmd/api-emulator/main.go
	GOOS=darwin GOARCH=amd64 go build -tags "-w -s" -o ./dist/tapi-emu-darwin ./terminal/cmd/api-emulator/main.go