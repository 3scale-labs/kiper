FROM golang:1.13
WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o 3scale-opa ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
USER 1001
COPY --from=0 /app/3scale-opa /app/3scale-opa
ENTRYPOINT ["/app/3scale-opa"]
