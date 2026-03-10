import { Request, Response, NextFunction } from 'express';
import { verifyAccessToken, JwtPayload } from '../lib/jwt';

export interface AuthRequest extends Request {
  user: JwtPayload;
}

export function requireAuth(req: Request, res: Response, next: NextFunction) {
  const header = req.headers.authorization;
  if (!header?.startsWith('Bearer ')) {
    res.status(401).json({ message: 'Unauthorized' });
    return;
  }

  try {
    const token = header.slice(7);
    (req as AuthRequest).user = verifyAccessToken(token);
    next();
  } catch {
    res.status(401).json({ message: 'Invalid or expired token' });
  }
}
