# Earthfile
VERSION 0.6

FROM golang:1.18.1-alpine

WORKDIR /app

all:
    BUILD +lint
    BUILD +test
    BUILD +coverage

code:
	# download deps
	COPY --dir go.mod go.sum ./
	RUN go mod download -x

	# copy in code
	COPY --dir ./ ./
	COPY --dir +protogen/gen ./

vendor:
	FROM +code

	RUN go mod vendor
	SAVE ARTIFACT /app /files


test:
	COPY +vendor/files ./

	WITH DOCKER --pull postgres --pull mysql
       RUN go test -mod=vendor -v ./... -coverprofile=coverage.out
    END

	SAVE ARTIFACT ./coverage.out AS LOCAL ./coverage.out

coverage:
	FROM +test

    ARG COMMIT_HASH=""
    ARG BRANCH_NAME=""
    ARG BUILD_NUMBER=""
    RUN curl -s https://codecov.io/bash > codecov.sh && chmod +x codecov.sh
    RUN --secret CODECOV_TOKEN=+secrets/CODECOV_TOKEN \
        ./codecov.sh -t "${CODECOV_TOKEN}" -B "${BRANCH_NAME}" -C "${COMMIT_HASH}" -b "${BUILD_NUMBER}"

lint:
    FROM +vendor

    COPY .golangci.yml ./

    # Runs golangci-lint with settings:
	RUN golangci-lint run -v

protogen:
    FROM +golang-base

    WORKDIR workspace

    # copy the proto files to generate
    COPY --dir testdata/ .
    COPY buf.work.yaml buf.gen.yaml .

    # generate the pbs
    RUN buf generate

    SAVE ARTIFACT /defs/gen / AS LOCAL ${OUT}

golang-base:

    WORKDIR /app
    ARG VERSION=dev

    # install gcc dependencies into alpine for CGO
    RUN apk add gcc musl-dev curl git openssh

    # install docker tools
    # https://docs.docker.com/engine/install/debian/
    RUN apk add --update --no-cache docker

    # install the go generator plugins
    RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    RUN export PATH="$PATH:$(go env GOPATH)/bin"

    # install buf from source
    RUN GO111MODULE=on GOBIN=/usr/local/bin go install github.com/bufbuild/buf/cmd/buf@v1.3.1

    # install linter
    # binary will be $(go env GOPATH)/bin/golangci-lint
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2
    RUN ls -la $(which golangci-lint)
