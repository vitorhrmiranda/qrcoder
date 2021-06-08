SERVERLESS_RUN = docker-compose run --rm serverless

setup:
	mkdir -p bin
	go build -o bin/qr
	$(SERVERLESS_RUN) sls deploy

generate:
	$(SERVERLESS_RUN) awslocal lambda invoke --function-name qrcoder-local-qrcoder --payload file://request.json response.json

decode:
	# make decode ext=svg
	cat response.json | jq .body | base64 -di > image.$(ext)
