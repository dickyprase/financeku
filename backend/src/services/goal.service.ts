import { goalRepository } from '../repositories';

export const goalService = {
  async create(userId: string, input: { name: string; target_amount: number; deadline?: string | null; tracking_mode: string; notes: string; wallet_ids: string[] }) {
    const goal = await goalRepository.create({
      user_id: userId,
      name: input.name,
      target_amount: input.target_amount,
      deadline: input.deadline || null,
      tracking_mode: input.tracking_mode as any,
      notes: input.notes || null,
    });

    for (const walletId of input.wallet_ids) {
      await goalRepository.linkWallet(goal.id, walletId);
    }

    return goal;
  },

  async list(userId: string) {
    return goalRepository.listByUser(userId);
  },

  async getById(userId: string, id: string) {
    return goalRepository.findById(id, userId);
  },

  async update(userId: string, id: string, input: { name: string; target_amount: number; deadline?: string | null; tracking_mode: string; notes: string; wallet_ids: string[] }) {
    const goal = await goalRepository.findById(id, userId);
    if (!goal) throw new Error('goal not found');

    const updated = await goalRepository.update(id, userId, {
      name: input.name,
      target_amount: input.target_amount,
      deadline: input.deadline || null,
      tracking_mode: input.tracking_mode as any,
      notes: input.notes || null,
      status: goal.status,
    });

    // Re-link wallets
    const existingWallets = await goalRepository.getGoalWallets(id);
    for (const gw of existingWallets) {
      await goalRepository.unlinkWallet(id, gw.wallet_id);
    }
    for (const walletId of input.wallet_ids) {
      await goalRepository.linkWallet(id, walletId);
    }

    return updated;
  },

  async delete(userId: string, id: string) {
    const goal = await goalRepository.findById(id, userId);
    if (!goal) throw new Error('goal not found');
    await goalRepository.delete(id, userId);
  },

  async getProgress(userId: string, id: string) {
    const goal = await goalRepository.findById(id, userId);
    if (!goal) throw new Error('goal not found');

    let currentAmount: number;
    if (goal.tracking_mode === 'manual') {
      currentAmount = parseFloat(goal.current_amount as any);
    } else {
      currentAmount = await goalRepository.calculateProgress(id, userId, goal.tracking_mode);
    }

    const targetAmount = parseFloat(goal.target_amount as any);
    let percentage = 0;
    if (targetAmount > 0) {
      percentage = Math.min((currentAmount / targetAmount) * 100, 100);
    }

    return {
      goal,
      current_amount: currentAmount,
      percentage: Math.round(percentage * 100) / 100,
    };
  },
};
