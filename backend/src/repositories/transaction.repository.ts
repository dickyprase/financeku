import { query } from '../config/database';
import { Transaction } from '../models';

export interface TransactionFilter {
  userId: string;
  walletId?: string;
  categoryId?: string;
  type?: string;
  dateFrom?: string;
  dateTo?: string;
  limit: number;
  offset: number;
}

export const transactionRepository = {
  async create(tx: Partial<Transaction>): Promise<Transaction> {
    const result = await query(
      `INSERT INTO transactions (user_id, wallet_id, category_id, type, amount, description, date, reference_id, reference_type)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
       RETURNING *`,
      [tx.user_id, tx.wallet_id, tx.category_id || null, tx.type, tx.amount, tx.description || null, tx.date, tx.reference_id || null, tx.reference_type || null]
    );
    return result.rows[0];
  },

  async findById(id: string, userId: string): Promise<Transaction | null> {
    const result = await query(
      'SELECT * FROM transactions WHERE id = $1 AND user_id = $2',
      [id, userId]
    );
    return result.rows[0] || null;
  },

  async list(filter: TransactionFilter): Promise<{ transactions: Transaction[]; total: number }> {
    let baseQuery = ' FROM transactions WHERE user_id = $1';
    const args: any[] = [filter.userId];
    let idx = 2;

    if (filter.walletId) {
      baseQuery += ` AND wallet_id = $${idx}`;
      args.push(filter.walletId);
      idx++;
    }
    if (filter.categoryId) {
      baseQuery += ` AND category_id = $${idx}`;
      args.push(filter.categoryId);
      idx++;
    }
    if (filter.type) {
      baseQuery += ` AND type = $${idx}`;
      args.push(filter.type);
      idx++;
    }
    if (filter.dateFrom) {
      baseQuery += ` AND date >= $${idx}`;
      args.push(filter.dateFrom);
      idx++;
    }
    if (filter.dateTo) {
      baseQuery += ` AND date <= $${idx}`;
      args.push(filter.dateTo);
      idx++;
    }

    const countResult = await query(`SELECT COUNT(*)${baseQuery}`, args);
    const total = parseInt(countResult.rows[0].count, 10);

    const selectQuery = `SELECT *${baseQuery} ORDER BY date DESC, created_at DESC LIMIT $${idx} OFFSET $${idx + 1}`;
    args.push(filter.limit, filter.offset);

    const result = await query(selectQuery, args);
    return { transactions: result.rows, total };
  },

  async update(id: string, userId: string, data: Partial<Transaction>): Promise<Transaction> {
    const result = await query(
      `UPDATE transactions SET wallet_id=$1, category_id=$2, type=$3, amount=$4, description=$5, date=$6, updated_at=NOW()
       WHERE id=$7 AND user_id=$8 RETURNING *`,
      [data.wallet_id, data.category_id || null, data.type, data.amount, data.description, data.date, id, userId]
    );
    return result.rows[0];
  },

  async delete(id: string, userId: string): Promise<void> {
    await query('DELETE FROM transactions WHERE id=$1 AND user_id=$2', [id, userId]);
  },

  async getMonthlySum(userId: string, type: string, yearMonth: string): Promise<number> {
    const result = await query(
      `SELECT COALESCE(SUM(amount), 0) as total FROM transactions
       WHERE user_id = $1 AND type = $2 AND TO_CHAR(date, 'YYYY-MM') = $3`,
      [userId, type, yearMonth]
    );
    return parseFloat(result.rows[0].total);
  },

  async getDailyExpense(userId: string, date: string): Promise<number> {
    const result = await query(
      `SELECT COALESCE(SUM(amount), 0) as total FROM transactions
       WHERE user_id = $1 AND type = 'expense' AND date = $2`,
      [userId, date]
    );
    return parseFloat(result.rows[0].total);
  },
};
