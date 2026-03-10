import amqp from 'amqplib';

export const SIGNUP_QUEUE = 'user.signup';
export const EMAIL_QUEUE = 'email.send';

let channel: amqp.Channel | null = null;

async function connect(): Promise<amqp.Channel> {
  const url = process.env.RABBITMQ_URL ?? 'amqp://guest:guest@localhost:5672';
  const conn = await amqp.connect(url);
  const ch = await conn.createChannel();
  await ch.assertQueue(SIGNUP_QUEUE, { durable: true });
  return ch;
}

export async function getPublishChannel(): Promise<amqp.Channel> {
  if (!channel) {
    channel = await connect();
  }
  return channel;
}

export async function publishEmail(payload: {
  to: string;
  subject: string;
  html: string;
}): Promise<void> {
  const ch = await getPublishChannel();
  await ch.assertQueue(EMAIL_QUEUE, { durable: true });
  ch.sendToQueue(EMAIL_QUEUE, Buffer.from(JSON.stringify(payload)), { persistent: true });
}

export async function publishSignup(payload: {
  jobId: string;
  email: string;
  password: string;
  name?: string;
}): Promise<void> {
  const ch = await getPublishChannel();
  ch.sendToQueue(SIGNUP_QUEUE, Buffer.from(JSON.stringify(payload)), { persistent: true });
}
