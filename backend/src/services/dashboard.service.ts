import { walletRepository, transactionRepository, overtimeRepository } from '../repositories';
import { dailyBudgetRepository } from '../repositories';

export const dashboardService = {
  async getSummary(userId: string) {
    const totalBalance = await walletRepository.getTotalBalance(userId);
    const currentMonth = new Date().toISOString().slice(0, 7); // YYYY-MM

    const monthlyIncome = await transactionRepository.getMonthlySum(userId, 'income', currentMonth);
    const monthlyExpense = await transactionRepository.getMonthlySum(userId, 'expense', currentMonth);
    const overtimePending = await overtimeRepository.getPendingTotal(userId);

    return {
      total_balance: totalBalance,
      monthly_income: monthlyIncome,
      monthly_expense: monthlyExpense,
      overtime_pending: overtimePending,
    };
  },

  async getCashflowReport(userId: string, month: string) {
    const income = await transactionRepository.getMonthlySum(userId, 'income', month);
    const expense = await transactionRepository.getMonthlySum(userId, 'expense', month);

    return {
      month,
      income,
      expense,
      net: income - expense,
    };
  },
};

export const dailyBudgetService = {
  async get(userId: string) {
    return dailyBudgetRepository.getByUser(userId);
  },

  async upsert(userId: string, input: { mode: string; manual_amount: number; formula_wallet_id: string; formula_days_remaining: number }) {
    return dailyBudgetRepository.upsert({
      user_id: userId,
      mode: input.mode as any,
      manual_amount: input.manual_amount || 0,
      formula_wallet_id: input.formula_wallet_id || null,
      formula_days_remaining: input.formula_days_remaining || 30,
    });
  },

  async getToday(userId: string) {
    const setting = await dailyBudgetRepository.getByUser(userId);

    if (!setting) {
      return { budget_amount: 0, spent_today: 0, remaining: 0 };
    }

    let budgetAmount = 0;

    if (setting.mode === 'manual') {
      budgetAmount = parseFloat(setting.manual_amount as any);
    } else if (setting.mode === 'formula' && setting.formula_wallet_id) {
      const wallet = await walletRepository.findById(setting.formula_wallet_id, userId);
      if (wallet) {
        const daysRemaining = setting.formula_days_remaining || 30;
        budgetAmount = parseFloat(wallet.balance as any) / daysRemaining;
      }
    }

    const today = new Date().toISOString().split('T')[0];
    const spentToday = await transactionRepository.getDailyExpense(userId, today);

    return {
      budget_amount: Math.round(budgetAmount * 100) / 100,
      spent_today: spentToday,
      remaining: Math.round((budgetAmount - spentToday) * 100) / 100,
    };
  },
};
