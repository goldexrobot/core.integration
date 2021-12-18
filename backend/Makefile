all: test swagger-callbacks apiv1
	@:

test:
	go test ./signature

swagger-callbacks:
	@MSYS_NO_PATHCONV=1 docker run --rm -v $(shell pwd)./../:/goldex:rw -it quay.io/goswagger/swagger generate spec -m -w /goldex/backend/callback -o /goldex/docs/swagger/callback/callback.swagger.json
	
apiv1:
	@$(MAKE) --no-print-directory -C api/v1/golang