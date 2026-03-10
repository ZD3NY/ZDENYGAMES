import 'dotenv/config';
import { PrismaClient } from '@prisma/client';
import bcrypt from 'bcrypt';
import crypto from 'crypto';

const prisma = new PrismaClient();

const PEPPER = process.env.PEPPER ?? 'change-me-pepper-in-production';

function applyPepper(password: string): string {
  return crypto.createHmac('sha256', PEPPER).update(password).digest('hex');
}

async function main() {
  const pepperedPassword = applyPepper('admin123');
  const password = await bcrypt.hash(pepperedPassword, 12);

  const user = await prisma.user.upsert({
    where: { email: 'admin@example.com' },
    update: { password },
    create: {
      email: 'admin@example.com',
      password,
      name: 'Admin',
    },
  });

  console.log('Seeded user:', user.email);
}

main()
  .catch(console.error)
  .finally(() => prisma.$disconnect());
