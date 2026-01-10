FROM golang:1.25.5-alpine3.23 AS build

WORKDIR /build
COPY . .

RUN adduser -D -H -g '' -u 10001 nonroot

RUN go mod download && \
    go mod verify && \
    GOOS=linux GOARCH=amd64 GOEXPERIMENT=jsonv2 GOEXPERIMENT=greenteagc go build -ldflags "-w -s -extldflags '-static -Wl,--strip-all,--gc-sections'" -o server


FROM scratch

COPY --from=build /build/server /server
COPY --from=build /etc/passwd /etc/passwd

USER nonroot

CMD ["/server"]
