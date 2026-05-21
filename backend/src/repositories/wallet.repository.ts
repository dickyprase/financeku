import { query } from '../config/database';
import { Wallet } from '../models';

export const walletRepository = {
  async create(wallet: Partial<Wallet>): Promise<Wallet> {
    const result = await query(
      `INSERT INTO wallets (user_id, name, balance, icon, color)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING *`,
      [wallet.user_id, wallet.name, wallet.balance || 0, wallet.icon || null, wallet.color || null]
    );
    return result.rows[0];
  },

  async findById(id: string, userId: string): Promise<Wallet | null> {
    const result = await query(
      'SELECT * FROM wallets WHERE id = $1 AND user_id = $2',
      [id, userId]
    );
    return result.rows[0] || null;
  },

  async listByUser(userId: string): Promise<Wallet[]> {
    const result = await query(
      'SELECT * FROM wallets WHERE user_id = $1 AND is_active = true ORDER BY created_at ASC',
      [userId]
    );
    return result.rows;
  },

  async update(id: string, userId: string, data: Partial<Wallet>): Promise<Wallet> {
    const result = await query(
      `UPDATE wallets SET name=$1, icon=$2, color=$3, updated_at=NOW()
       WHERE id=$4 AND user_id=$5 RETURNING *`,
      [data.name, data.icon, data.color, id, userId]
    );
    return result.rows[0];
  },

  async updateBalance(id: string, amount: number): Promise<void> {
    await query(
      'UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE id = $2',
      [amount, id]
    );
  },

  async setBalance(id: string, balance: number): Promise<void> {
    await query(
      'UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2',
      [balance, id]
    );
  },

  async delete(id: string, userId: string): Promise<void> {
    await query(
      'UPDATE wallets SET is_active = false, updated_at = NOW() WHERE id = $1 AND user_id = $2',
      [id, userId]
    );
  },

  async getTotalBalance(userId: string): Promise<number> {
    const result = await query(
      'SELECT COALESCE(SUM(balance), 0) as total FROM wallets WHERE user_id = $1 AND is_active = true',
      [userId]
    );
    return parseFloat(result.rows[0].total);
  },
};
