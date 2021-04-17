FROM golang:1.16.3-alpine3.13 as builder
WORKDIR /go/src/account-rest-service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/main cmd/account-rest-service/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/bin/main /go/bin/main
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/go/bin/main"]
