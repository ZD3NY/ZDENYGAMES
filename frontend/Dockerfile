FROM node:22-alpine AS builder

RUN npm install -g @quasar/cli

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .
RUN quasar build

FROM nginx:stable-alpine

COPY --from=builder /app/dist/spa /usr/share/nginx/html
COPY nginx-main.conf /etc/nginx/nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf.template

EXPOSE 80

CMD ["/bin/sh", "-c", "envsubst '${BACKEND_URL} ${BACKEND_HOST} ${TETRIS_URL} ${TETRIS_HOST}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]
