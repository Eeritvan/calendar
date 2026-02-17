ARG GO_VERSION=1.26.0
ARG ALPINE_VERSION=3.23
ARG BUN_VERSION=1.3.8

# build frontend
FROM oven/bun:${BUN_VERSION}-alpine AS frontend-build

WORKDIR /app

COPY frontend/package.json frontend/bun.lock ./

RUN bun i --frozen-lockfile

COPY frontend .

RUN bun run build


# build backend
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS backend-build

WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify

COPY backend .
COPY --from=frontend-build /app/dist ./dist

RUN adduser -D -H -g '' -u 10001 nonroot

RUN GOOS=linux GOARCH=amd64 GOEXPERIMENT=jsonv2 GOEXPERIMENT=greenteagc go build -ldflags "-w -s -extldflags '-static -Wl,--strip-all,--gc-sections'" -o server


# final stage
FROM scratch

COPY --from=backend-build /build/server /server
COPY --from=backend-build /etc/passwd /etc/passwd

USER nonroot

CMD ["/server"]
