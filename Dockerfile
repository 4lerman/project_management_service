FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache make

RUN make build

EXPOSE 8080

CMD ["sh", "-c", "while ! nc -z db 5432; do sleep 1; done; make migrate-up && make run"]
