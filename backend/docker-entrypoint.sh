#!/bin/sh
set -e

echo "Running database migrations..."
cd /app
pnpm --filter backend exec prisma migrate deploy

echo "Starting server..."
exec node /app/backend/dist/index.js
