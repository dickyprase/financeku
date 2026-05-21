import { query } from '../config/database';
import { User } from '../models';

export const userRepository = {
  async create(user: Partial<User>): Promise<User> {
    const result = await query(
      `INSERT INTO users (name, email, password, role, phone, telegram, salary, meal_allowance)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
       RETURNING *`,
      [user.name, user.email, user.password, user.role || 'user', user.phone || null, user.telegram || null, user.salary || 0, user.meal_allowance || 0]
    );
    return result.rows[0];
  },

  async findByEmail(email: string): Promise<User | null> {
    const result = await query('SELECT * FROM users WHERE email = $1', [email]);
    return result.rows[0] || null;
  },

  async findById(id: string): Promise<User | null> {
    const result = await query('SELECT * FROM users WHERE id = $1', [id]);
    return result.rows[0] || null;
  },

  async update(id: string, data: Partial<User>): Promise<User> {
    const result = await query(
      `UPDATE users SET name=$1, phone=$2, telegram=$3, salary=$4, meal_allowance=$5, updated_at=NOW()
       WHERE id=$6 RETURNING *`,
      [data.name, data.phone, data.telegram, data.salary, data.meal_allowance, id]
    );
    return result.rows[0];
  },

  async updatePassword(id: string, hashedPassword: string): Promise<void> {
    await query('UPDATE users SET password=$1, updated_at=NOW() WHERE id=$2', [hashedPassword, id]);
  },

  async list(limit: number, offset: number): Promise<{ users: User[]; total: number }> {
    const countResult = await query('SELECT COUNT(*) FROM users');
    const total = parseInt(countResult.rows[0].count, 10);

    const result = await query(
      `SELECT id, name, email, role, phone, telegram, salary, meal_allowance, is_active, created_at, updated_at
       FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
      [limit, offset]
    );
    return { users: result.rows, total };
  },

  async setActive(id: string, active: boolean): Promise<void> {
    await query('UPDATE users SET is_active=$1, updated_at=NOW() WHERE id=$2', [active, id]);
  },

  async delete(id: string): Promise<void> {
    await query('DELETE FROM users WHERE id=$1', [id]);
  },
};
