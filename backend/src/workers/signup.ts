import amqp from 'amqplib';
import bcrypt from 'bcrypt';
import { prisma } from '../lib/prisma';
import { applyPepper } from '../lib/crypto';
import { SIGNUP_QUEUE, publishEmail } from '../lib/rabbitmq';

const SALT_ROUNDS = 12;

export async function startSignupWorker(): Promise<void> {
  const url = process.env.RABBITMQ_URL ?? 'amqp://guest:guest@localhost:5672';

  let conn: amqp.ChannelModel | null = null;
  for (let attempt = 1; attempt <= 15; attempt++) {
    try {
      conn = await amqp.connect(url);
      console.log('Signup worker connected to RabbitMQ');
      break;
    } catch {
      console.log(`RabbitMQ not ready, retrying (${attempt}/15)...`);
      await new Promise((r) => setTimeout(r, 3000));
    }
  }

  if (!conn) {
    throw new Error('Could not connect to RabbitMQ after retries');
  }

  const ch = await conn.createChannel();
  await ch.assertQueue(SIGNUP_QUEUE, { durable: true });
  ch.prefetch(1);

  console.log('Signup worker listening on queue:', SIGNUP_QUEUE);

  ch.consume(SIGNUP_QUEUE, async (msg: amqp.Message | null) => {
    if (!msg) return;

    const raw = msg.content.toString();
    let jobId: string | undefined;

    try {
      const { jobId: id, email, password, name } = JSON.parse(raw) as {
        jobId: string;
        email: string;
        password: string;
        name?: string;
      };
      jobId = id;

      const peppered = applyPepper(password);
      const hashed = await bcrypt.hash(peppered, SALT_ROUNDS);

      const user = await prisma.user.create({
        data: { email, password: hashed, name: name ?? null },
      });

      await prisma.signupJob.update({
        where: { id: jobId },
        data: { status: 'done', userId: user.id },
      });

      await publishEmail({
        to: email,
        subject: 'Welcome!',
        html: `
          <h2>Welcome${user.name ? `, ${user.name}` : ''}!</h2>
          <p>Your account has been created successfully.</p>
          <p>You can now <a href="${process.env.FRONTEND_URL ?? 'http://localhost'}">sign in</a>.</p>
        `,
      });
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      console.error('Signup worker error:', message);
      if (jobId) {
        await prisma.signupJob
          .update({ where: { id: jobId }, data: { status: 'failed', error: message } })
          .catch(() => {});
      }
    }

    ch.ack(msg);
  });
}
