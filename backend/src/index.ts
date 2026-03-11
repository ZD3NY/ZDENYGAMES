import 'dotenv/config';
import express from 'express';
import cors from 'cors';
import { router } from './routes';
import { startSignupWorker } from './workers/signup';
import { getPublishChannel } from './lib/rabbitmq';

const app = express();
const PORT = process.env.PORT ?? 3000;

app.use(cors({ origin: process.env.FRONTEND_URL ?? 'http://localhost:9000' }));
app.use(express.json());

app.use('/api', router);

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});

startSignupWorker().catch((err) => {
  console.error('Failed to start signup worker:', err);
});

getPublishChannel().catch((err) => {
  console.error('Failed to initialize publish channel:', err);
});
