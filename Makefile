all: backend terminal
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: backend

backend:
	@$(MAKE) --no-print-directory -C backend

terminal:
	@$(MAKE) --no-print-directory -C terminal