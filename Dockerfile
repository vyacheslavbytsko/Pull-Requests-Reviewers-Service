FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o server ./cmd/server/main.go
EXPOSE 8080
CMD ["/app/server"]