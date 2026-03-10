import { Router, IRouter } from 'express';
import { authRouter } from './auth';
import { scoresRouter } from './scores';

export const router: IRouter = Router();

router.get('/health', (_req, res) => {
  res.json({ status: 'ok' });
});

router.use('/auth', authRouter);
router.use('/scores', scoresRouter);
