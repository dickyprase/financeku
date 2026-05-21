import { query } from '../config/database';
import { OvertimeRecord } from '../models';

export const overtimeRepository = {
  async create(record: Partial<OvertimeRecord>): Promise<OvertimeRecord> {
    const result = await query(
      `INSERT INTO overtime_records (user_id, date, hours, is_holiday, amount, meal_amount, total_amount, period_start, period_end, notes)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
       RETURNING *`,
      [record.user_id, record.date, record.hours, record.is_holiday, record.amount, record.meal_amount, record.total_amount, record.period_start, record.period_end, record.notes || null]
    );
    return result.rows[0];
  },

  async findById(id: string, userId: string): Promise<OvertimeRecord | null> {
    const result = await query(
      'SELECT * FROM overtime_records WHERE id = $1 AND user_id = $2',
      [id, userId]
    );
    return result.rows[0] || null;
  },

  async listByUser(userId: string, month: string, limit: number, offset: number): Promise<{ records: OvertimeRecord[]; total: number }> {
    let baseQuery = ' FROM overtime_records WHERE user_id = $1';
    const args: any[] = [userId];
    let idx = 2;

    if (month) {
      baseQuery += ` AND TO_CHAR(date, 'YYYY-MM') = $${idx}`;
      args.push(month);
      idx++;
    }

    const countResult = await query(`SELECT COUNT(*)${baseQuery}`, args);
    const total = parseInt(countResult.rows[0].count, 10);

    const selectQuery = `SELECT *${baseQuery} ORDER BY date DESC LIMIT $${idx} OFFSET $${idx + 1}`;
    args.push(limit, offset);

    const result = await query(selectQuery, args);
    return { records: result.rows, total };
  },

  async update(record: Partial<OvertimeRecord> & { id: string; user_id: string }): Promise<OvertimeRecord> {
    const result = await query(
      `UPDATE overtime_records SET date=$1, hours=$2, is_holiday=$3, amount=$4, meal_amount=$5,
       total_amount=$6, period_start=$7, period_end=$8, notes=$9, updated_at=NOW()
       WHERE id=$10 AND user_id=$11 AND is_disbursed=false RETURNING *`,
      [record.date, record.hours, record.is_holiday, record.amount, record.meal_amount, record.total_amount, record.period_start, record.period_end, record.notes, record.id, record.user_id]
    );
    return result.rows[0];
  },

  async delete(id: string, userId: string): Promise<void> {
    await query('DELETE FROM overtime_records WHERE id=$1 AND user_id=$2 AND is_disbursed=false', [id, userId]);
  },

  async getPeriodRecords(userId: string, periodStart: string, periodEnd: string): Promise<OvertimeRecord[]> {
    const result = await query(
      `SELECT * FROM overtime_records
       WHERE user_id = $1 AND period_start = $2 AND period_end = $3 AND is_disbursed = false
       ORDER BY date ASC`,
      [userId, periodStart, periodEnd]
    );
    return result.rows;
  },

  async disbursePeriod(userId: string, periodStart: string, periodEnd: string): Promise<void> {
    await query(
      `UPDATE overtime_records SET is_disbursed = true, disbursed_at = NOW(), updated_at = NOW()
       WHERE user_id = $1 AND period_start = $2 AND period_end = $3 AND is_disbursed = false`,
      [userId, periodStart, periodEnd]
    );
  },

  async getPendingTotal(userId: string): Promise<number> {
    const result = await query(
      'SELECT COALESCE(SUM(total_amount), 0) as total FROM overtime_records WHERE user_id = $1 AND is_disbursed = false',
      [userId]
    );
    return parseFloat(result.rows[0].total);
  },
};
