import { query } from '../config/database';
import { Goal, GoalWallet } from '../models';

export const goalRepository = {
  async create(goal: Partial<Goal>): Promise<Goal> {
    const result = await query(
      `INSERT INTO goals (user_id, name, target_amount, deadline, tracking_mode, notes)
       VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING *`,
      [goal.user_id, goal.name, goal.target_amount, goal.deadline || null, goal.tracking_mode || 'manual', goal.notes || null]
    );
    return result.rows[0];
  },

  async findById(id: string, userId: string): Promise<Goal | null> {
    const result = await query(
      'SELECT * FROM goals WHERE id = $1 AND user_id = $2',
      [id, userId]
    );
    return result.rows[0] || null;
  },

  async listByUser(userId: string): Promise<Goal[]> {
    const result = await query(
      'SELECT * FROM goals WHERE user_id = $1 ORDER BY created_at DESC',
      [userId]
    );
    return result.rows;
  },

  async update(id: string, userId: string, data: Partial<Goal>): Promise<Goal> {
    const result = await query(
      `UPDATE goals SET name=$1, target_amount=$2, deadline=$3, tracking_mode=$4, notes=$5, status=$6, updated_at=NOW()
       WHERE id=$7 AND user_id=$8 RETURNING *`,
      [data.name, data.target_amount, data.deadline || null, data.tracking_mode, data.notes, data.status || 'active', id, userId]
    );
    return result.rows[0];
  },

  async updateCurrentAmount(id: string, amount: number): Promise<void> {
    await query('UPDATE goals SET current_amount=$1, updated_at=NOW() WHERE id=$2', [amount, id]);
  },

  async delete(id: string, userId: string): Promise<void> {
    await query('DELETE FROM goals WHERE id=$1 AND user_id=$2', [id, userId]);
  },

  async linkWallet(goalId: string, walletId: string): Promise<void> {
    await query(
      'INSERT INTO goal_wallets (goal_id, wallet_id) VALUES ($1, $2) ON CONFLICT DO NOTHING',
      [goalId, walletId]
    );
  },

  async unlinkWallet(goalId: string, walletId: string): Promise<void> {
    await query('DELETE FROM goal_wallets WHERE goal_id=$1 AND wallet_id=$2', [goalId, walletId]);
  },

  async getGoalWallets(goalId: string): Promise<GoalWallet[]> {
    const result = await query(
      'SELECT * FROM goal_wallets WHERE goal_id = $1',
      [goalId]
    );
    return result.rows;
  },

  async calculateProgress(goalId: string, userId: string, trackingMode: string): Promise<number> {
    if (trackingMode === 'all_wallet') {
      const result = await query(
        'SELECT COALESCE(SUM(balance), 0) as total FROM wallets WHERE user_id = $1 AND is_active = true',
        [userId]
      );
      return parseFloat(result.rows[0].total);
    }

    // single_wallet or multiple_wallet
    const result = await query(
      `SELECT COALESCE(SUM(w.balance), 0) as total
       FROM wallets w
       INNER JOIN goal_wallets gw ON gw.wallet_id = w.id
       WHERE gw.goal_id = $1 AND w.is_active = true`,
      [goalId]
    );
    return parseFloat(result.rows[0].total);
  },
};
