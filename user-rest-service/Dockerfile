FROM golang:1.15.2-alpine3.12 as builder
WORKDIR /go/src/user-rest-service
RUN addgroup -g 10001 -S appgroup \
    && adduser -u 10001 -G appgroup -S appuser
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/main cmd/user-rest-service/main.go

FROM scratch
COPY --from=builder /go/bin/main /go/bin/main
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
EXPOSE 8080
USER appuser
ENTRYPOINT ["/go/bin/main"]
