FROM golang:1.26.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION
ARG BUILD_TIME
RUN go build \
    -ldflags "-X DevOpsMiniProject/util.Version=${VERSION} -X DevOpsMiniProject/util.buildTime=${BUILD_TIME} -X DevOpsMiniProject/util.CommitSHA=${COMMIT_SHA}" \
    -o /app/main ./cmd

EXPOSE 3010
CMD ["./main"]