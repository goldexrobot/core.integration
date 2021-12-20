all: swagger-api
	@:

swagger-api:
	$(info Swagger for API)
	@MSYS_NO_PATHCONV=1 docker run --rm -v $(shell pwd)./../:/goldex:rw -it quay.io/goswagger/swagger generate spec -m -w /goldex/terminal/api/v1 -o /goldex/docs/swagger/terminal/api/v1/api.swagger.json
