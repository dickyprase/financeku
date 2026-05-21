import { incomeRepository, walletRepository, transactionRepository } from '../repositories';

export const incomeService = {
  async create(userId: string, input: { wallet_id: string; amount: number; source: string; description: string; date: string }) {
    const walletId = input.wallet_id || null;

    if (walletId) {
      const wallet = await walletRepository.findById(walletId, userId);
      if (!wallet) throw new Error('wallet not found');
    }

    const income = await incomeRepository.create({
      user_id: userId,
      wallet_id: walletId,
      amount: input.amount,
      source: input.source,
      description: input.description || null,
      date: input.date,
    });

    // If wallet linked, create transaction and update balance
    if (walletId) {
      await transactionRepository.create({
        user_id: userId,
        wallet_id: walletId,
        type: 'income',
        amount: input.amount,
        description: `Income: ${input.source}`,
        date: input.date,
        reference_id: income.id,
        reference_type: 'income',
      });
      await walletRepository.updateBalance(walletId, input.amount);
    }

    return income;
  },

  async list(userId: string, limit: number, offset: number) {
    return incomeRepository.listByUser(userId, limit, offset);
  },

  async getById(userId: string, id: string) {
    return incomeRepository.findById(id, userId);
  },

  async delete(userId: string, id: string) {
    const income = await incomeRepository.findById(id, userId);
    if (!income) throw new Error('income not found');

    // Reverse wallet balance if linked
    if (income.wallet_id) {
      await walletRepository.updateBalance(income.wallet_id, -parseFloat(income.amount as any));
    }

    await incomeRepository.delete(id, userId);
  },
};
