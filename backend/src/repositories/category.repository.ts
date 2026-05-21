import { query } from '../config/database';
import { Category } from '../models';

export const categoryRepository = {
  async create(cat: Partial<Category>): Promise<Category> {
    const result = await query(
      `INSERT INTO categories (user_id, name, type, icon, color, budget_limit)
       VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING *`,
      [cat.user_id, cat.name, cat.type, cat.icon || null, cat.color || null, cat.budget_limit || 0]
    );
    return result.rows[0];
  },

  async findById(id: string): Promise<Category | null> {
    const result = await query('SELECT * FROM categories WHERE id = $1', [id]);
    return result.rows[0] || null;
  },

  async listByUser(userId: string, type?: string): Promise<Category[]> {
    if (type) {
      const result = await query(
        `SELECT * FROM categories WHERE (user_id = $1 OR is_default = true) AND type = $2
         ORDER BY is_default DESC, name ASC`,
        [userId, type]
      );
      return result.rows;
    }
    const result = await query(
      `SELECT * FROM categories WHERE user_id = $1 OR is_default = true
       ORDER BY type, is_default DESC, name ASC`,
      [userId]
    );
    return result.rows;
  },

  async update(id: string, userId: string, data: Partial<Category>): Promise<Category> {
    const result = await query(
      `UPDATE categories SET name=$1, icon=$2, color=$3, budget_limit=$4, updated_at=NOW()
       WHERE id=$5 AND user_id=$6 AND is_default=false RETURNING *`,
      [data.name, data.icon, data.color, data.budget_limit, id, userId]
    );
    return result.rows[0];
  },

  async delete(id: string, userId: string): Promise<void> {
    await query('DELETE FROM categories WHERE id=$1 AND user_id=$2 AND is_default=false', [id, userId]);
  },
};
