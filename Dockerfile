FROM golang:1.26

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o crawler ./cmd/crawler
RUN go build -o searchapi ./cmd/searchapi

CMD ["./searchapi"]