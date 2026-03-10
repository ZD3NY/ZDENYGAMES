import 'dotenv/config';
import amqp from 'amqplib';
import nodemailer from 'nodemailer';

const QUEUE = 'email.send';

const transporter = nodemailer.createTransport({
  host: process.env.SMTP_HOST ?? 'mailhog',
  port: Number(process.env.SMTP_PORT ?? 1025),
  secure: false,
  auth:
    process.env.SMTP_USER && process.env.SMTP_PASS
      ? { user: process.env.SMTP_USER, pass: process.env.SMTP_PASS }
      : undefined,
});

interface EmailJob {
  to: string;
  subject: string;
  html: string;
}

async function start() {
  const url = process.env.RABBITMQ_URL ?? 'amqp://guest:guest@localhost:5672';

  let conn: amqp.ChannelModel | null = null;
  for (let attempt = 1; attempt <= 15; attempt++) {
    try {
      conn = await amqp.connect(url);
      console.log('Email service connected to RabbitMQ');
      break;
    } catch {
      console.log(`RabbitMQ not ready, retrying (${attempt}/15)...`);
      await new Promise((r) => setTimeout(r, 3000));
    }
  }

  if (!conn) throw new Error('Could not connect to RabbitMQ');

  const ch = await conn.createChannel();
  await ch.assertQueue(QUEUE, { durable: true });
  ch.prefetch(1);

  console.log('Email service listening on queue:', QUEUE);

  ch.consume(QUEUE, async (msg: amqp.Message | null) => {
    if (!msg) return;

    try {
      const job = JSON.parse(msg.content.toString()) as EmailJob;

      await transporter.sendMail({
        from: process.env.SMTP_FROM ?? 'noreply@business.local',
        to: job.to,
        subject: job.subject,
        html: job.html,
      });

      console.log(`Email sent to ${job.to}: ${job.subject}`);
    } catch (err) {
      console.error('Failed to send email:', err);
    }

    ch.ack(msg);
  });
}

start().catch((err) => {
  console.error('Email service failed to start:', err);
  process.exit(1);
});
