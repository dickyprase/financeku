import { query } from '../config/database';
import { DailyBudgetSetting, ActivityLog, SiteSetting } from '../models';

export const dailyBudgetRepository = {
  async getByUser(userId: string): Promise<DailyBudgetSetting | null> {
    const result = await query(
      'SELECT * FROM daily_budget_settings WHERE user_id = $1',
      [userId]
    );
    return result.rows[0] || null;
  },

  async upsert(setting: Partial<DailyBudgetSetting>): Promise<DailyBudgetSetting> {
    const result = await query(
      `INSERT INTO daily_budget_settings (user_id, mode, manual_amount, formula_wallet_id, formula_days_remaining)
       VALUES ($1, $2, $3, $4, $5)
       ON CONFLICT (user_id) DO UPDATE SET
         mode = EXCLUDED.mode,
         manual_amount = EXCLUDED.manual_amount,
         formula_wallet_id = EXCLUDED.formula_wallet_id,
         formula_days_remaining = EXCLUDED.formula_days_remaining,
         updated_at = NOW()
       RETURNING *`,
      [setting.user_id, setting.mode, setting.manual_amount || 0, setting.formula_wallet_id || null, setting.formula_days_remaining || 30]
    );
    return result.rows[0];
  },
};

export const activityLogRepository = {
  async create(log: Partial<ActivityLog>): Promise<ActivityLog> {
    const result = await query(
      `INSERT INTO activity_logs (user_id, action, entity_type, entity_id, details, ip_address)
       VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING *`,
      [log.user_id || null, log.action, log.entity_type || null, log.entity_id || null, log.details || null, log.ip_address || null]
    );
    return result.rows[0];
  },

  async list(limit: number, offset: number): Promise<{ logs: ActivityLog[]; total: number }> {
    const countResult = await query('SELECT COUNT(*) FROM activity_logs');
    const total = parseInt(countResult.rows[0].count, 10);

    const result = await query(
      'SELECT * FROM activity_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2',
      [limit, offset]
    );
    return { logs: result.rows, total };
  },
};

export const siteSettingRepository = {
  async getAll(): Promise<SiteSetting[]> {
    const result = await query('SELECT * FROM site_settings ORDER BY key');
    return result.rows;
  },

  async getByKey(key: string): Promise<SiteSetting | null> {
    const result = await query('SELECT * FROM site_settings WHERE key = $1', [key]);
    return result.rows[0] || null;
  },

  async update(key: string, value: string): Promise<void> {
    await query('UPDATE site_settings SET value = $1, updated_at = NOW() WHERE key = $2', [value, key]);
  },
};
