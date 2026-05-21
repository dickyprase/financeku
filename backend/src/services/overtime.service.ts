import { overtimeRepository, userRepository, incomeRepository, walletRepository, transactionRepository } from '../repositories';
import { calculateOvertime, calculatePeriod } from '../utils/overtime-calc';

export const overtimeService = {
  async calculate(userId: string, hours: number, isHoliday: boolean) {
    const user = await userRepository.findById(userId);
    if (!user) throw new Error('user not found');

    const result = calculateOvertime(hours, user.salary, user.meal_allowance, isHoliday);
    return {
      hours,
      is_holiday: isHoliday,
      base_amount: result.baseAmount,
      meal_amount: result.mealAmount,
      total_amount: result.totalAmount,
    };
  },

  async create(userId: string, input: { date: string; hours: number; is_holiday: boolean; notes: string }) {
    const user = await userRepository.findById(userId);
    if (!user) throw new Error('user not found');

    const result = calculateOvertime(input.hours, user.salary, user.meal_allowance, input.is_holiday);
    const { periodStart, periodEnd } = calculatePeriod(input.date);

    const record = await overtimeRepository.create({
      user_id: userId,
      date: input.date,
      hours: input.hours,
      is_holiday: input.is_holiday,
      amount: result.baseAmount,
      meal_amount: result.mealAmount,
      total_amount: result.totalAmount,
      period_start: periodStart,
      period_end: periodEnd,
      notes: input.notes,
    });

    return record;
  },

  async update(userId: string, id: string, input: { date: string; hours: number; is_holiday: boolean; notes: string }) {
    const record = await overtimeRepository.findById(id, userId);
    if (!record) throw new Error('overtime record not found');
    if (record.is_disbursed) throw new Error('cannot update disbursed overtime record');

    const user = await userRepository.findById(userId);
    if (!user) throw new Error('user not found');

    const result = calculateOvertime(input.hours, user.salary, user.meal_allowance, input.is_holiday);
    const { periodStart, periodEnd } = calculatePeriod(input.date);

    const updated = await overtimeRepository.update({
      id,
      user_id: userId,
      date: input.date,
      hours: input.hours,
      is_holiday: input.is_holiday,
      amount: result.baseAmount,
      meal_amount: result.mealAmount,
      total_amount: result.totalAmount,
      period_start: periodStart,
      period_end: periodEnd,
      notes: input.notes,
    });

    return updated;
  },

  async delete(userId: string, id: string) {
    const record = await overtimeRepository.findById(id, userId);
    if (!record) throw new Error('overtime record not found');
    if (record.is_disbursed) throw new Error('cannot delete disbursed overtime record');
    await overtimeRepository.delete(id, userId);
  },

  async list(userId: string, month: string, limit: number, offset: number) {
    return overtimeRepository.listByUser(userId, month, limit, offset);
  },

  async getById(userId: string, id: string) {
    return overtimeRepository.findById(id, userId);
  },

  async disburse(userId: string, input: { period_start: string; period_end: string; wallet_id: string }) {
    const records = await overtimeRepository.getPeriodRecords(userId, input.period_start, input.period_end);
    if (records.length === 0) throw new Error('no overtime records found for this period');

    const totalAmount = records.reduce((sum, r) => sum + parseFloat(r.total_amount as any), 0);

    // Mark as disbursed
    await overtimeRepository.disbursePeriod(userId, input.period_start, input.period_end);

    // Create income record
    const income = await incomeRepository.create({
      user_id: userId,
      wallet_id: input.wallet_id,
      amount: totalAmount,
      source: 'Overtime',
      description: `Overtime disbursement period ${input.period_start} to ${input.period_end}`,
      date: new Date().toISOString().split('T')[0],
      is_from_overtime: true,
      overtime_period_start: input.period_start,
      overtime_period_end: input.period_end,
    });

    // Create transaction and update wallet
    if (input.wallet_id) {
      await transactionRepository.create({
        user_id: userId,
        wallet_id: input.wallet_id,
        type: 'income',
        amount: totalAmount,
        description: `Overtime disbursement: ${input.period_start} - ${input.period_end}`,
        date: new Date().toISOString().split('T')[0],
        reference_id: income.id,
        reference_type: 'overtime_disbursement',
      });

      await walletRepository.updateBalance(input.wallet_id, totalAmount);
    }
  },
};
