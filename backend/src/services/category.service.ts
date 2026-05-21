import { categoryRepository } from '../repositories';

export const categoryService = {
  async create(userId: string, input: { name: string; type: string; icon: string; color: string; budget_limit: number }) {
    return categoryRepository.create({
      user_id: userId,
      name: input.name,
      type: input.type as any,
      icon: input.icon || null,
      color: input.color || null,
      budget_limit: input.budget_limit || 0,
    });
  },

  async list(userId: string, type?: string) {
    return categoryRepository.listByUser(userId, type);
  },

  async getById(id: string) {
    return categoryRepository.findById(id);
  },

  async update(userId: string, id: string, input: { name: string; icon: string; color: string; budget_limit: number }) {
    const cat = await categoryRepository.findById(id);
    if (!cat) throw new Error('category not found');
    if (cat.is_default) throw new Error('cannot update default category');

    return categoryRepository.update(id, userId, {
      name: input.name,
      icon: input.icon || null,
      color: input.color || null,
      budget_limit: input.budget_limit || 0,
    });
  },

  async delete(userId: string, id: string) {
    const cat = await categoryRepository.findById(id);
    if (!cat) throw new Error('category not found');
    if (cat.is_default) throw new Error('cannot delete default category');
    await categoryRepository.delete(id, userId);
  },
};
