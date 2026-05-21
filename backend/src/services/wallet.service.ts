import { walletRepository, transactionRepository } from '../repositories';

export const walletService = {
  async create(userId: string, input: { name: string; balance: number; icon: string; color: string }) {
    return walletRepository.create({
      user_id: userId,
      name: input.name,
      balance: input.balance,
      icon: input.icon || null,
      color: input.color || null,
    });
  },

  async list(userId: string) {
    return walletRepository.listByUser(userId);
  },

  async getById(userId: string, id: string) {
    return walletRepository.findById(id, userId);
  },

  async update(userId: string, id: string, input: { name: string; icon: string; color: string }) {
    const wallet = await walletRepository.findById(id, userId);
    if (!wallet) throw new Error('wallet not found');

    return walletRepository.update(id, userId, {
      name: input.name,
      icon: input.icon || null,
      color: input.color || null,
    });
  },

  async delete(userId: string, id: string) {
    const wallet = await walletRepository.findById(id, userId);
    if (!wallet) throw new Error('wallet not found');
    await walletRepository.delete(id, userId);
  },

  async transfer(userId: string, input: { from_wallet_id: string; to_wallet_id: string; amount: number; admin_fee: number; description: string; date: string }) {
    const fromWallet = await walletRepository.findById(input.from_wallet_id, userId);
    if (!fromWallet) throw new Error('source wallet not found');

    const toWallet = await walletRepository.findById(input.to_wallet_id, userId);
    if (!toWallet) throw new Error('destination wallet not found');

    const totalDebit = input.amount + input.admin_fee;
    if (parseFloat(fromWallet.balance as any) < totalDebit) {
      throw new Error('insufficient balance in source wallet');
    }

    // Debit source
    await walletRepository.updateBalance(input.from_wallet_id, -totalDebit);

    // Credit destination
    await walletRepository.updateBalance(input.to_wallet_id, input.amount);

    // Create debit transaction
    const debitTx = await transactionRepository.create({
      user_id: userId,
      wallet_id: input.from_wallet_id,
      type: 'transfer',
      amount: totalDebit,
      description: `Transfer to wallet: ${input.description}`,
      date: input.date,
    });

    // Create credit transaction
    await transactionRepository.create({
      user_id: userId,
      wallet_id: input.to_wallet_id,
      type: 'transfer',
      amount: input.amount,
      description: `Transfer from wallet: ${input.description}`,
      date: input.date,
      reference_id: debitTx.id,
    });
  },

  async getTotalBalance(userId: string) {
    return walletRepository.getTotalBalance(userId);
  },
};
