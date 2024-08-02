FROM golang:1.22-alpine

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o /bin/app

EXPOSE 3000

ENTRYPOINT ["/bin/app"]

CMD ["./app"]
