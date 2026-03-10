import crypto from 'crypto';

const PEPPER = process.env.PEPPER ?? 'change-me-pepper-in-production';

export function applyPepper(password: string): string {
  return crypto.createHmac('sha256', PEPPER).update(password).digest('hex');
}
