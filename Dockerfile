FROM golang:alpine3.19

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@v1.49.0

COPY . .
RUN go mod tidy

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]