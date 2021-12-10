all: test callback-spec apiv1
	@:

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

test:
	go test ./signature

callback-spec:
	@MSYS_NO_PATHCONV=1 docker run --rm -v $(shell pwd):/goldex:rw -it quay.io/goswagger/swagger generate spec -m -w /goldex/callback -o /goldex/docs/swagger/callback/callback.swagger.json
	
apiv1:
	@$(MAKE) --no-print-directory -C api/v1/golang