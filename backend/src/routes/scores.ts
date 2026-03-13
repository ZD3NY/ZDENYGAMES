import { Router, IRouter, Request, Response } from 'express';
import { prisma } from '../lib/prisma';
import { requireAuth, AuthRequest } from '../middleware/auth';

export const scoresRouter: IRouter = Router();

scoresRouter.post('/', requireAuth, async (req: Request, res: Response) => {
  const { score, lines } = req.body as { score?: number; lines?: number };
  const userId = (req as AuthRequest).user.sub;

  if (typeof score !== 'number' || typeof lines !== 'number') {
    res.status(400).json({ message: 'score and lines are required numbers' });
    return;
  }

  try {
    const existing = await prisma.score.findUnique({ where: { userId } });
    if (existing && existing.score >= score) {
      res.status(200).json(existing);
      return;
    }
    const entry = await prisma.score.upsert({
      where: { userId },
      update: { score, lines, createdAt: new Date() },
      create: { userId, score, lines },
    });
    res.status(201).json(entry);
  } catch (err) {
    console.error('Score submit error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

scoresRouter.get('/leaderboard', async (_req: Request, res: Response) => {
  try {
    const scores = await prisma.score.findMany({
      orderBy: { score: 'desc' },
      take: 10,
      include: { user: { select: { name: true, email: true } } },
    });

    const leaderboard = scores.map((s, i) => ({
      rank: i + 1,
      name: s.user.name ?? s.user.email.split('@')[0],
      score: s.score,
      lines: s.lines,
      date: s.createdAt,
    }));

    res.json(leaderboard);
  } catch (err) {
    console.error('Leaderboard error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});
