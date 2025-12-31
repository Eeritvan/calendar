FROM oven/bun:1.3.5-alpine AS base

# prepare turborepo
FROM base AS prepare
WORKDIR /app

COPY . .

RUN bunx turbo prune web-client --docker


# building
FROM base AS build
WORKDIR /app

COPY --from=prepare /app/out/json/ /temp/prod
RUN --mount=type=cache,target=/root/.bun \
    cd /temp/prod && bun i -p --frozen-lockfile

COPY --from=prepare /app/out/json/ .
RUN --mount=type=cache,target=/root/.bun \
    bun i --frozen-lockfile

COPY --from=prepare /app/out/full/ .
RUN bun run build

# run
FROM base AS runner
WORKDIR /app

COPY --from=build /temp/prod/node_modules node_modules
COPY --from=build /app/apps/web/.output .output
COPY --from=build /app/apps/web/package.json .

ENTRYPOINT ["bun", "run", "start"]

ENV PORT=3000
CMD ["--port", "$PORT"]
