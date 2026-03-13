import { Router, IRouter, Request, Response } from 'express';
import { prisma } from '../lib/prisma';
import { requireAuth, AuthRequest } from '../middleware/auth';

export const scoresRouter: IRouter = Router();

// POST /api/scores/tetris
scoresRouter.post('/tetris', requireAuth, async (req: Request, res: Response) => {
  const { score, lines } = req.body as { score?: number; lines?: number };
  const userId = (req as AuthRequest).user.sub;

  if (typeof score !== 'number' || typeof lines !== 'number') {
    res.status(400).json({ message: 'score and lines are required numbers' });
    return;
  }

  try {
    const existing = await prisma.scoreTetris.findUnique({ where: { userId } });
    if (existing && existing.score >= score) {
      res.status(200).json(existing);
      return;
    }
    const entry = await prisma.scoreTetris.upsert({
      where: { userId },
      update: { score, lines, createdAt: new Date() },
      create: { userId, score, lines },
    });
    res.status(201).json(entry);
  } catch (err) {
    console.error('Tetris score submit error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

// POST /api/scores/wolfpack
scoresRouter.post('/wolfpack', requireAuth, async (req: Request, res: Response) => {
  const { score, waves } = req.body as { score?: number; waves?: number };
  const userId = (req as AuthRequest).user.sub;

  if (typeof score !== 'number' || typeof waves !== 'number') {
    res.status(400).json({ message: 'score and waves are required numbers' });
    return;
  }

  try {
    const existing = await prisma.scoreWolfpack.findUnique({ where: { userId } });
    if (existing && existing.score >= score) {
      res.status(200).json(existing);
      return;
    }
    const entry = await prisma.scoreWolfpack.upsert({
      where: { userId },
      update: { score, waves, createdAt: new Date() },
      create: { userId, score, waves },
    });
    res.status(201).json(entry);
  } catch (err) {
    console.error('Wolfpack score submit error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

// GET /api/scores/tetris/leaderboard
scoresRouter.get('/tetris/leaderboard', async (_req: Request, res: Response) => {
  try {
    const scores = await prisma.scoreTetris.findMany({
      orderBy: { score: 'desc' },
      take: 10,
      include: { user: { select: { name: true, email: true } } },
    });

    res.json(scores.map((s, i) => ({
      rank: i + 1,
      name: s.user.name ?? s.user.email.split('@')[0],
      score: s.score,
      lines: s.lines,
      date: s.createdAt,
    })));
  } catch (err) {
    console.error('Tetris leaderboard error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});

// GET /api/scores/wolfpack/leaderboard
scoresRouter.get('/wolfpack/leaderboard', async (_req: Request, res: Response) => {
  try {
    const scores = await prisma.scoreWolfpack.findMany({
      orderBy: { score: 'desc' },
      take: 10,
      include: { user: { select: { name: true, email: true } } },
    });

    res.json(scores.map((s, i) => ({
      rank: i + 1,
      name: s.user.name ?? s.user.email.split('@')[0],
      score: s.score,
      waves: s.waves,
      date: s.createdAt,
    })));
  } catch (err) {
    console.error('Wolfpack leaderboard error:', err);
    res.status(500).json({ message: 'Internal server error' });
  }
});
