FROM golang:alpine
LABEL authors="markjforte"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /opnsense-dns-register

CMD ["/opnsense-dns-register"]