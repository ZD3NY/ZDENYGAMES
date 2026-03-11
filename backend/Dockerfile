FROM node:22-alpine

# native deps for bcrypt + openssl for Prisma
RUN apk add --no-cache python3 make g++ openssl

WORKDIR /app

COPY package.json ./
COPY prisma ./prisma

RUN npm install

COPY src ./src
COPY tsconfig.json ./

RUN npx prisma generate && npx tsc

EXPOSE 3000

CMD ["/bin/sh", "-c", "npx prisma migrate deploy && node dist/index.js"]
