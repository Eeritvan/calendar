FROM golang:1.25.5-alpine3.22 AS build

WORKDIR /build
COPY . .

RUN adduser -D -H -g '' -u 10001 nonroot

RUN go mod download && \
    go mod verify && \
    GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -extldflags '-static -Wl,--strip-all,--gc-sections'" -o calendar


FROM scratch

COPY --from=build /build/calendar /calendar
COPY --from=build /etc/passwd /etc/passwd

USER nonroot

CMD ["/calendar"]
