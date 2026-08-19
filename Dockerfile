FROM golang:1.26.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION
ARG BUILD_TIME
ARG COMMIT_SHA
RUN CGO_ENABLED=0 go build \
    -ldflags "-X DevOpsMiniProject/util.Version=${VERSION} -X DevOpsMiniProject/util.buildTime=${BUILD_TIME} -X DevOpsMiniProject/util.CommitSHA=${COMMIT_SHA}" \
    -o /app/main ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3010
CMD ["./main"]