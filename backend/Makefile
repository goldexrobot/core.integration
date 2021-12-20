all: test callback golang-apiv1
	@:

test:
	$(info Go test)
	@go test ./signature

.PHONY: callback
callback:
	$(info Swagger for callbacks)
	@MSYS_NO_PATHCONV=1 docker run --rm -v $(shell pwd)./../:/goldex:rw -it quay.io/goswagger/swagger generate spec -m -w /goldex/backend/callback -o /goldex/docs/swagger/backend/callback.swagger.json
	
golang-apiv1:
	$(info Golang API v1 from proto)
	@$(MAKE) --no-print-directory -C api/v1/golang

