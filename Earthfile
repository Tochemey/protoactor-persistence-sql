# Earthfile

FROM golang:1.15-alpine3.13

WORKDIR /app

code:
	# add git to known hosts
	RUN mkdir -p /root/.ssh && \
		chmod 700 /root/.ssh && \
		ssh-keyscan github.com >> /root/.ssh/known_hosts

	RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

	# download deps
	COPY -dir go.mod go.sum ./
	RUN --ssh go mod download -x

	# copy in code
	COPY -dir ./ ./

vendor:
	FROM +code

	# install dependencies (using host machine ssh context)
	RUN --ssh go mod vendor
	SAVE ARTIFACT /app /files