import { pool } from './config/database';
import * as fs from 'fs';
import * as path from 'path';

async function migrate() {
  console.log('Running migrations...');

  const migrationFile = path.join(__dirname, '..', 'migrations', '001_initial_schema.sql');
  const migrationSQL = fs.readFileSync(migrationFile, 'utf-8');

  try {
    await pool.query(migrationSQL);
    console.log('Migrations completed successfully!');
  } catch (err) {
    console.error('Migration failed:', err);
    process.exit(1);
  }

  // Check for --seed flag
  if (process.argv.includes('--seed')) {
    console.log('Running seed...');
    const seedFile = path.join(__dirname, '..', 'migrations', '002_seed_data.sql');
    const seedSQL = fs.readFileSync(seedFile, 'utf-8');

    try {
      await pool.query(seedSQL);
      console.log('Seed completed successfully!');
    } catch (err) {
      console.error('Seed failed:', err);
      process.exit(1);
    }
  }

  await pool.end();
}

migrate();
