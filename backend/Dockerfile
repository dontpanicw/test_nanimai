FROM golang:1.23.4-alpine as build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache make

RUN go build -o bin/backend ./backend

FROM alpine:3.19 AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /build/bin/backend /app/backend

CMD ["/app/backend", "--rest-addr", ":8080", "--grpc-addr", ":9090"]