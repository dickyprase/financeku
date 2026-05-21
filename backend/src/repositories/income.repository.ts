import { query } from '../config/database';
import { Income } from '../models';

export const incomeRepository = {
  async create(income: Partial<Income>): Promise<Income> {
    const result = await query(
      `INSERT INTO incomes (user_id, wallet_id, amount, source, description, date, is_from_overtime, overtime_period_start, overtime_period_end)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
       RETURNING *`,
      [income.user_id, income.wallet_id || null, income.amount, income.source, income.description || null, income.date, income.is_from_overtime || false, income.overtime_period_start || null, income.overtime_period_end || null]
    );
    return result.rows[0];
  },

  async findById(id: string, userId: string): Promise<Income | null> {
    const result = await query(
      'SELECT * FROM incomes WHERE id = $1 AND user_id = $2',
      [id, userId]
    );
    return result.rows[0] || null;
  },

  async listByUser(userId: string, limit: number, offset: number): Promise<{ incomes: Income[]; total: number }> {
    const countResult = await query('SELECT COUNT(*) FROM incomes WHERE user_id = $1', [userId]);
    const total = parseInt(countResult.rows[0].count, 10);

    const result = await query(
      'SELECT * FROM incomes WHERE user_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3',
      [userId, limit, offset]
    );
    return { incomes: result.rows, total };
  },

  async update(id: string, userId: string, data: Partial<Income>): Promise<Income> {
    const result = await query(
      `UPDATE incomes SET wallet_id=$1, amount=$2, source=$3, description=$4, date=$5, updated_at=NOW()
       WHERE id=$6 AND user_id=$7 AND is_from_overtime=false RETURNING *`,
      [data.wallet_id, data.amount, data.source, data.description, data.date, id, userId]
    );
    return result.rows[0];
  },

  async delete(id: string, userId: string): Promise<void> {
    await query('DELETE FROM incomes WHERE id=$1 AND user_id=$2', [id, userId]);
  },
};
