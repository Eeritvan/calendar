FROM golang:1.24.3-alpine3.21 AS build

WORKDIR /build
COPY . .

RUN go mod download && \
    go mod verify && \
    GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -extldflags '-static -Wl,--strip-all,--gc-sections'" -o calendar-backend


FROM scratch

COPY --from=build /build/calendar-backend /calendar-backend

CMD ["/calendar-backend"]
