FROM golang:1.24.3

RUN apt update && apt install -y build-essential

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o app ./cmd/run_benchmark

CMD ["go", "run", "./cmd/run_benchmark"]
