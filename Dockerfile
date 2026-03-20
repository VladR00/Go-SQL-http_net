FROM golang:1.26.1-alpine3.23 AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o ./build/executable/app ./cmd/main.go

FROM alpine:3.23

WORKDIR /app

COPY --from=build /app/build/executable/app ./app
COPY --from=build /app/.env ./.env

EXPOSE 8080

CMD ["/app/app"]