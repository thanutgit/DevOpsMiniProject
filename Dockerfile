FROM golang:1.26.6 AS builder
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

FROM alpine:3.24
RUN apk add --no-cache ca-certificates
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /app/main .
USER 10001
EXPOSE 3010
CMD ["./main"]