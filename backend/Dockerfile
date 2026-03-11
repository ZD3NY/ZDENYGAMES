FROM node:22-alpine

# native deps for bcrypt + openssl for Prisma
RUN apk add --no-cache python3 make g++ openssl
RUN npm install -g pnpm

WORKDIR /app

# copy workspace manifests
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY backend/package.json ./backend/
COPY frontend/package.json ./frontend/

# copy prisma schema before install so @prisma/client postinstall can generate
COPY backend/prisma ./backend/prisma

RUN pnpm install --frozen-lockfile --filter backend...

COPY backend/src ./backend/src
COPY backend/tsconfig.json ./backend/tsconfig.json

RUN pnpm --filter backend exec prisma generate && pnpm --filter backend build

EXPOSE 3000

CMD ["/bin/sh", "-c", "cd /app && pnpm --filter backend exec prisma migrate deploy && exec node /app/backend/dist/index.js"]
