FROM golang:1.17-alpine AS build

WORKDIR /app

COPY src/go.mod .
COPY src/go.sum .

RUN go mod download

COPY src .

RUN go build -o ./build/rsoi-lab01 ./main.go

FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/build/rsoi-lab01 .
EXPOSE 8080
CMD ["./rsoi-lab01"]
