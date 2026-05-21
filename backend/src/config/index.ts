import dotenv from 'dotenv';

dotenv.config();

export const config = {
  server: {
    port: parseInt(process.env.SERVER_PORT || '8080', 10),
    env: process.env.APP_ENV || 'development',
  },
  db: {
    host: process.env.DB_HOST || 'localhost',
    port: parseInt(process.env.DB_PORT || '5432', 10),
    user: process.env.DB_USER || 'postgres',
    password: process.env.DB_PASSWORD || 'postgres',
    name: process.env.DB_NAME || 'financeku',
    sslMode: process.env.DB_SSLMODE || 'disable',
  },
  jwt: {
    secret: process.env.JWT_SECRET || 'change-me-in-production',
    accessExpMinutes: parseInt(process.env.JWT_ACCESS_EXP_MINUTES || '15', 10),
    refreshExpDays: parseInt(process.env.JWT_REFRESH_EXP_DAYS || '7', 10),
  },
  cors: {
    allowedOrigins: process.env.ALLOWED_ORIGINS || '*',
  },
};
