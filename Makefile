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