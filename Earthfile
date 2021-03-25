# Earthfile

FROM golang:1.15.10-buster
WORKDIR /app

all:
    BUILD +lint
    BUILD +test
    BUILD +coverage

code:
	# download deps
	COPY -dir go.mod go.sum ./
	RUN go mod download -x

	# copy in code
	COPY -dir ./ ./

vendor:
	FROM +code

	# install dependencies (using host machine ssh context)
	RUN go mod vendor
	SAVE ARTIFACT /app /files

test:
	COPY +vendor/files ./
	RUN go test -mod=vendor -v ./... -coverprofile=coverage.out
	SAVE IMAGE
	SAVE ARTIFACT ./coverage.out AS LOCAL ./coverage.out

coverage:
	FROM +test

	ARG COMMIT_HASH=""
	ARG BRANCH_NAME=""
	ARG BUILD_NUMBER=""
	ARG CODECOV_TOKEN=""
	RUN curl -s https://codecov.io/bash > ./codecov.sh && chmod +x ./codecov.sh
	RUN ./codecov.sh -t "${CODECOV_TOKEN}" -B "${BRANCH_NAME}" -C "${COMMIT_HASH}" -b "${BUILD_NUMBER}"

lint:
    FROM +vendor

     # binary will be $(go env GOPATH)/bin/golangci-lint
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.38.0

    # Runs golangci-lint with settings:
	RUN golangci-lint run -v
