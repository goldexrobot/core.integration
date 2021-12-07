all: test apiv1
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %


test:
	go test ./signature

apiv1:
	@$(MAKE) --no-print-directory -C api/v1/golang