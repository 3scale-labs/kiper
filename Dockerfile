FROM golang:1.13
WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kiper ./

FROM alpine:latest
RUN apk --no-cache add ca-certificates
USER 1001
COPY --from=0 /app/kiper /app/kiper
ENTRYPOINT ["/app/kiper"]
