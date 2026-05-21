import { transactionRepository, walletRepository } from '../repositories';
import { TransactionFilter } from '../repositories/transaction.repository';

export const transactionService = {
  async create(userId: string, input: { wallet_id: string; category_id: string; type: string; amount: number; description: string; date: string }) {
    const wallet = await walletRepository.findById(input.wallet_id, userId);
    if (!wallet) throw new Error('wallet not found');

    const tx = await transactionRepository.create({
      user_id: userId,
      wallet_id: input.wallet_id,
      category_id: input.category_id || null,
      type: input.type as any,
      amount: input.amount,
      description: input.description || null,
      date: input.date,
    });

    // Update wallet balance
    if (input.type === 'income') {
      await walletRepository.updateBalance(input.wallet_id, input.amount);
    } else if (input.type === 'expense') {
      await walletRepository.updateBalance(input.wallet_id, -input.amount);
    }

    return tx;
  },

  async list(userId: string, filter: Omit<TransactionFilter, 'userId'>) {
    return transactionRepository.list({ ...filter, userId });
  },

  async getById(userId: string, id: string) {
    return transactionRepository.findById(id, userId);
  },

  async update(userId: string, id: string, input: { wallet_id: string; category_id: string; type: string; amount: number; description: string; date: string }) {
    const existing = await transactionRepository.findById(id, userId);
    if (!existing) throw new Error('transaction not found');

    // Reverse old balance
    if (existing.type === 'income') {
      await walletRepository.updateBalance(existing.wallet_id, -parseFloat(existing.amount as any));
    } else if (existing.type === 'expense') {
      await walletRepository.updateBalance(existing.wallet_id, parseFloat(existing.amount as any));
    }

    const updated = await transactionRepository.update(id, userId, {
      wallet_id: input.wallet_id,
      category_id: input.category_id || null,
      type: input.type as any,
      amount: input.amount,
      description: input.description || null,
      date: input.date,
    });

    // Apply new balance
    if (input.type === 'income') {
      await walletRepository.updateBalance(input.wallet_id, input.amount);
    } else if (input.type === 'expense') {
      await walletRepository.updateBalance(input.wallet_id, -input.amount);
    }

    return updated;
  },

  async delete(userId: string, id: string) {
    const existing = await transactionRepository.findById(id, userId);
    if (!existing) throw new Error('transaction not found');

    // Reverse balance
    if (existing.type === 'income') {
      await walletRepository.updateBalance(existing.wallet_id, -parseFloat(existing.amount as any));
    } else if (existing.type === 'expense') {
      await walletRepository.updateBalance(existing.wallet_id, parseFloat(existing.amount as any));
    }

    await transactionRepository.delete(id, userId);
  },
};
