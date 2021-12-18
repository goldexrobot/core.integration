all: backend
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: backend

backend:
	@$(MAKE) --no-print-directory -C backend
