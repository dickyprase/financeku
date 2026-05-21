import { userRepository } from '../repositories';
import { siteSettingRepository, activityLogRepository } from '../repositories';
import { hashPassword, checkPassword } from '../utils/hash';

export const profileService = {
  async update(userId: string, input: { name: string; phone: string; telegram: string; salary: number; meal_allowance: number }) {
    const user = await userRepository.findById(userId);
    if (!user) throw new Error('user not found');

    return userRepository.update(userId, {
      name: input.name,
      phone: input.phone || null,
      telegram: input.telegram || null,
      salary: input.salary || 0,
      meal_allowance: input.meal_allowance || 0,
    });
  },

  async changePassword(userId: string, oldPassword: string, newPassword: string) {
    const user = await userRepository.findById(userId);
    if (!user) throw new Error('user not found');

    const valid = await checkPassword(oldPassword, user.password);
    if (!valid) throw new Error('current password is incorrect');

    const hashed = await hashPassword(newPassword);
    await userRepository.updatePassword(userId, hashed);
  },
};

export const adminService = {
  async listUsers(limit: number, offset: number) {
    return userRepository.list(limit, offset);
  },

  async createUser(input: { name: string; email: string; password: string; role: string; salary: number; meal_allowance: number }) {
    const existing = await userRepository.findByEmail(input.email);
    if (existing) throw new Error('email already registered');

    const hashed = await hashPassword(input.password);
    return userRepository.create({
      name: input.name,
      email: input.email,
      password: hashed,
      role: input.role || 'user',
      salary: input.salary || 0,
      meal_allowance: input.meal_allowance || 0,
    });
  },

  async updateUser(id: string, input: { name: string; role: string; salary: number; meal_allowance: number }) {
    const user = await userRepository.findById(id);
    if (!user) throw new Error('user not found');

    return userRepository.update(id, {
      name: input.name,
      phone: user.phone,
      telegram: user.telegram,
      salary: input.salary || 0,
      meal_allowance: input.meal_allowance || 0,
    });
  },

  async deleteUser(id: string) {
    await userRepository.delete(id);
  },

  async resetPassword(id: string, newPassword: string) {
    const hashed = await hashPassword(newPassword);
    await userRepository.updatePassword(id, hashed);
  },

  async getSettings() {
    return siteSettingRepository.getAll();
  },

  async updateSettings(settings: Record<string, string>) {
    for (const [key, value] of Object.entries(settings)) {
      await siteSettingRepository.update(key, value);
    }
  },

  async getActivityLogs(limit: number, offset: number) {
    return activityLogRepository.list(limit, offset);
  },
};
