# syntax=docker/dockerfile:1

FROM node:22-alpine AS build
WORKDIR /app/frontend

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY --from=api-docs openapi.yaml /app/api-docs/openapi.yaml
COPY . .
RUN npm run build

FROM nginxinc/nginx-unprivileged:stable-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/frontend/dist /usr/share/nginx/html
EXPOSE 8080
