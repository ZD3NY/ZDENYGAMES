import { Router, IRouter, Request, Response } from 'express';
import bcrypt from 'bcrypt';
import { prisma } from '../lib/prisma';
import { signAccessToken, signRefreshToken, verifyRefreshToken } from '../lib/jwt';
import { applyPepper } from '../lib/crypto';
import { requireAuth, AuthRequest } from '../middleware/auth';
import { publishSignup } from '../lib/rabbitmq';

export const authRouter: IRouter = Router();

const SALT_ROUNDS = 12;
const REFRESH_TOKEN_EXPIRES_DAYS = 7;

function refreshTokenExpiresAt(): Date {
  const date = new Date();
  date.setDate(date.getDate() + REFRESH_TOKEN_EXPIRES_DAYS);
  return date;
}

authRouter.post('/sign-up', async (req: Request, res: Response) => {
  const { email, password, name } = req.body as {
    email?: string;
    password?: string;
    name?: string;
  };

  if (!email || !password) {
    res.status(400).json({ message: 'Email and password are required' });
    return;
  }

  try {
    const existing = await prisma.user.findUnique({ where: { email } });
    if (existing) {
      res.status(409).json({ message: 'Email already registered' });
      return;
    }

    const job = await prisma.signupJob.create({ data: { email } });
    await publishSignup({ jobId: job.id, email, password, name });

    res.status(202).json({ jobId: job.id });
  } catch (err) {
    console.error('Sign-up error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

authRouter.get('/sign-up/status/:jobId', async (req: Request<{ jobId: string }>, res: Response) => {
  const { jobId } = req.params;

  try {
    const job = await prisma.signupJob.findUnique({ where: { id: jobId } });
    if (!job) {
      res.status(404).json({ message: 'Job not found' });
      return;
    }

    if (job.status === 'pending') {
      res.json({ status: 'pending' });
      return;
    }

    if (job.status === 'failed') {
      res.json({ status: 'failed', error: job.error });
      return;
    }

    // done — issue tokens and return session
    const user = await prisma.user.findUnique({ where: { id: job.userId! } });
    if (!user) {
      res.status(500).json({ message: 'User not found' });
      return;
    }

    const payload = { sub: user.id, email: user.email };
    const accessToken = signAccessToken(payload);
    const refreshToken = signRefreshToken(payload);

    await prisma.refreshToken.create({
      data: { token: refreshToken, userId: user.id, expiresAt: refreshTokenExpiresAt() },
    });

    res.json({
      status: 'done',
      accessToken,
      refreshToken,
      user: { id: user.id, email: user.email, name: user.name },
    });
  } catch (err) {
    console.error('Sign-up status error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

authRouter.post('/sign-in', async (req: Request, res: Response) => {
  const { email, password } = req.body as { email?: string; password?: string };

  if (!email || !password) {
    res.status(400).json({ message: 'Email and password are required' });
    return;
  }

  try {
    const user = await prisma.user.findUnique({ where: { email } });
    if (!user) {
      res.status(401).json({ message: 'Invalid email or password' });
      return;
    }

    // Pepper is applied before bcrypt compare; bcrypt stores the salt internally
    const pepperedPassword = applyPepper(password);
    const valid = await bcrypt.compare(pepperedPassword, user.password);
    if (!valid) {
      res.status(401).json({ message: 'Invalid email or password' });
      return;
    }

    const payload = { sub: user.id, email: user.email };
    const accessToken = signAccessToken(payload);
    const refreshToken = signRefreshToken(payload);

    await prisma.refreshToken.create({
      data: { token: refreshToken, userId: user.id, expiresAt: refreshTokenExpiresAt() },
    });

    res.json({
      accessToken,
      refreshToken,
      user: { id: user.id, email: user.email, name: user.name },
    });
  } catch (err) {
    console.error('Sign-in error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

authRouter.post('/refresh', async (req: Request, res: Response) => {
  const { refreshToken } = req.body as { refreshToken?: string };

  if (!refreshToken) {
    res.status(400).json({ message: 'Refresh token is required' });
    return;
  }

  try {
    const payload = verifyRefreshToken(refreshToken);

    const stored = await prisma.refreshToken.findUnique({ where: { token: refreshToken } });
    if (!stored || stored.expiresAt < new Date()) {
      if (stored) await prisma.refreshToken.delete({ where: { id: stored.id } });
      res.status(401).json({ message: 'Invalid or expired refresh token' });
      return;
    }

    // Rotate: delete old refresh token and issue new pair
    await prisma.refreshToken.delete({ where: { id: stored.id } });

    const newPayload = { sub: payload.sub, email: payload.email };
    const newAccessToken = signAccessToken(newPayload);
    const newRefreshToken = signRefreshToken(newPayload);

    await prisma.refreshToken.create({
      data: { token: newRefreshToken, userId: payload.sub, expiresAt: refreshTokenExpiresAt() },
    });

    res.json({ accessToken: newAccessToken, refreshToken: newRefreshToken });
  } catch {
    res.status(401).json({ message: 'Invalid or expired refresh token' });
  }
});

authRouter.post('/sign-out', requireAuth, async (req: Request, res: Response) => {
  const { refreshToken } = req.body as { refreshToken?: string };

  if (refreshToken) {
    await prisma.refreshToken.deleteMany({ where: { token: refreshToken } }).catch(() => {});
  }

  // Optionally revoke all sessions for this user
  const authReq = req as AuthRequest;
  if (!refreshToken && authReq.user?.sub) {
    await prisma.refreshToken.deleteMany({ where: { userId: authReq.user.sub } }).catch(() => {});
  }

  res.json({ message: 'Signed out successfully' });
});
