FROM node:22-alpine AS builder

RUN npm install -g pnpm

WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/
COPY backend/package.json ./backend/

RUN pnpm install --frozen-lockfile --filter frontend...

COPY frontend ./frontend

RUN pnpm --filter frontend build

FROM nginx:alpine

COPY --from=builder /app/frontend/dist/spa /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
