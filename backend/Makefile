all: test callback apiv1
	@:

test:
	go test ./signature

.PHONY: callback
callback:
	@MSYS_NO_PATHCONV=1 docker run --rm -v $(shell pwd)./../:/goldex:rw -it quay.io/goswagger/swagger generate spec -m -w /goldex/backend/callback -o /goldex/docs/swagger/backend/callback.swagger.json
	
apiv1:
	@$(MAKE) --no-print-directory -C api/v1/golang