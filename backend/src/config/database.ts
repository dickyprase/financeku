import { Pool } from 'pg';
import { config } from './index';

export const pool = new Pool({
  host: config.db.host,
  port: config.db.port,
  user: config.db.user,
  password: config.db.password,
  database: config.db.name,
  ssl: config.db.sslMode === 'require' ? { rejectUnauthorized: false } : false,
  max: 25,
  idleTimeoutMillis: 60000,
  connectionTimeoutMillis: 5000,
});

pool.on('connect', () => {
  console.log('Database connected successfully');
});

pool.on('error', (err) => {
  console.error('Unexpected database error:', err);
});

export const query = (text: string, params?: any[]) => pool.query(text, params);
