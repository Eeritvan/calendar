FROM oven/bun:alpine AS build

WORKDIR /app

COPY package.json bun.lock ./
RUN bun i --frozen-lockfile

COPY . .

ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=${VITE_BACKEND_URL}

ARG VITE_WEBSOCKET_URL
ENV VITE_WEBSOCKET_URL=${VITE_WEBSOCKET_URL}

RUN bun run build


FROM oven/bun:alpine

WORKDIR /app

COPY --from=build /app/build ./build
COPY package.json bun.lock ./

RUN bun i -p --frozen-lockfile && \
    rm -rf /root/.bun/install/cache /app/bun.lock /tmp/*

CMD ["bun", "run", "start"]