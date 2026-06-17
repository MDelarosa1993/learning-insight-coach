FROM golang:1.22-bullseye AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /bin/coach main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /bin/coach /bin/coach
COPY evals/ evals/
COPY data/ data/

EXPOSE 8080

CMD ["/bin/coach"]