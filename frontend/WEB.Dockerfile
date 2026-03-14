FROM docker.io/oven/bun:latest AS build

WORKDIR /app

COPY . .

ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN bun install
RUN bun run build-only

FROM docker.io/library/nginx:1.29.5-alpine AS base

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
EXPOSE 443
