# Stage 1: compile game to WebAssembly
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=js GOARCH=wasm go build -o web/main.wasm .

# Copy wasm_exec.js from the Go toolchain
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/

# Stage 2: serve with nginx
FROM nginx:stable-alpine

COPY --from=builder /app/web /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080
