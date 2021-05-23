# Earthfile

FROM golang:1.16.2-buster

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

	RUN curl -s https://codecov.io/bash > ./codecov.sh && chmod +x ./codecov.sh
	RUN ./codecov.sh -B "${BRANCH_NAME}" -C "${COMMIT_HASH}" -b "${BUILD_NUMBER}"

lint:
    FROM +vendor

     # binary will be $(go env GOPATH)/bin/golangci-lint
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.38.0

    # Runs golangci-lint with settings:
	RUN golangci-lint run -v

protogen:
	FROM namely/protoc-all:1.34_0

	ARG OUT="gen"
    ENV OUT=${OUT}

	COPY --dir protos /defs

	ARG GENERATE = entrypoint.sh --no-google-includes -l go -o ./${OUT}

	RUN ${GENERATE} -i protos -d protos

    SAVE ARTIFACT /defs/gen / AS LOCAL ${OUT}
