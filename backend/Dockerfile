FROM golang:1.23.4-alpine as build

WORKDIR /build

COPY ../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN apk add --no-cache make

RUN go build -o bin/backend ./backend

FROM alpine AS runner

WORKDIR backend

RUN apk add --no-cache curl

COPY --from=build /build/bin/backend /build/backend

CMD ["/build/backend"]