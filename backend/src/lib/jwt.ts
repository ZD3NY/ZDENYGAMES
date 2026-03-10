import jwt from 'jsonwebtoken';

const ACCESS_SECRET = process.env.JWT_SECRET ?? 'change-me-in-production';
const REFRESH_SECRET = process.env.JWT_REFRESH_SECRET ?? 'change-me-refresh-in-production';

// Access tokens are short-lived; refresh tokens are long-lived
const ACCESS_EXPIRES_IN = '15m';
const REFRESH_EXPIRES_IN = '7d';

export interface JwtPayload {
  sub: string;
  email: string;
}

export function signAccessToken(payload: JwtPayload): string {
  return jwt.sign(payload, ACCESS_SECRET, { expiresIn: ACCESS_EXPIRES_IN });
}

export function signRefreshToken(payload: JwtPayload): string {
  return jwt.sign(payload, REFRESH_SECRET, { expiresIn: REFRESH_EXPIRES_IN });
}

export function verifyAccessToken(token: string): JwtPayload {
  return jwt.verify(token, ACCESS_SECRET) as JwtPayload;
}

export function verifyRefreshToken(token: string): JwtPayload {
  return jwt.verify(token, REFRESH_SECRET) as JwtPayload;
}
