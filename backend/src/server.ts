import app from './app';
import { config } from './config';
import { pool } from './config/database';

async function start() {
  // Test database connection
  try {
    await pool.query('SELECT 1');
    console.log('Database connected successfully');
  } catch (err) {
    console.error('Failed to connect to database:', err);
    process.exit(1);
  }

  // Start server
  app.listen(config.server.port, () => {
    console.log(`FinanceKu API server starting on :${config.server.port}`);
    console.log(`Environment: ${config.server.env}`);
  });
}

start();
