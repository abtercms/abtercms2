PACKAGES := $(shell find . -name *.go | grep -v -e vendor | xargs -n1 dirname | sort -u)

default: build

.PHONY: generate
generate:
	go generate $(PACKAGES)

.PHONY: build
build: generate
	golangci-lint run
	go test ./...
	sam build

.PHONY: aws-list-tables
aws-list-tables:
	aws dynamodb list-tables --endpoint-url http://localhost:8000

.PHONY: aws-create-table-websites
aws-create-table-websites:
	aws dynamodb create-table --table-name websites --attribute-definitions AttributeName=pk,AttributeType=S --key-schema AttributeName=pk,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://localhost:8000

.PHONY: curl-list-websites
curl-list-websites:
	curl http:/127.0.0.1:3000/websites

.PHONY: curl-create-website
curl-create-website:
	curl -d '{"foo":"bar"}' -H "Content-Type: application/json" -X POST http:/127.0.0.1:3000/websites

.PHONY: curl-update-website-abc
curl-update-website-abc:
	curl -d '{"pk":"abc", "foo":"bar"}' -H "Content-Type: application/json" -X PUT http:/127.0.0.1:3000/websites/abc

.PHONY: curl-get-website-abc
curl-get-website-abc:
	curl http:/127.0.0.1:3000/websites/abc

.PHONY: sam-local
sam-local: build
	sam local start-api

.PHONY: local-dynamodb
local-dynamodb:
	java -Djava.library.path=./.local/DynamoDBLocal_lib -jar ./.local/DynamoDBLocal.jar -sharedDb

.PHONY: clean
clean:
	rm -rvf pkg/mocks websites/mocks .aws-sam

