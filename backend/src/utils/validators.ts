import { z } from 'zod';

export const registerSchema = z.object({
  name: z.string().min(1, 'name is required'),
  email: z.string().email('email must be a valid email address'),
  password: z.string().min(6, 'password must be at least 6 characters'),
});

export const loginSchema = z.object({
  email: z.string().min(1, 'email is required'),
  password: z.string().min(1, 'password is required'),
});

export const refreshSchema = z.object({
  refresh_token: z.string().min(1, 'refresh_token is required'),
});

export const walletSchema = z.object({
  name: z.string().min(1, 'name is required'),
  balance: z.number().optional().default(0),
  icon: z.string().optional().default(''),
  color: z.string().optional().default(''),
});

export const transferSchema = z.object({
  from_wallet_id: z.string().uuid('from_wallet_id is required'),
  to_wallet_id: z.string().uuid('to_wallet_id is required'),
  amount: z.number().min(1, 'amount must be at least 1'),
  admin_fee: z.number().optional().default(0),
  description: z.string().optional().default(''),
  date: z.string().min(1, 'date is required'),
});

export const categorySchema = z.object({
  name: z.string().min(1, 'name is required'),
  type: z.enum(['income', 'expense'], { message: 'type must be income or expense' }),
  icon: z.string().optional().default(''),
  color: z.string().optional().default(''),
  budget_limit: z.number().optional().default(0),
});

export const transactionSchema = z.object({
  wallet_id: z.string().uuid('wallet_id is required'),
  category_id: z.string().optional().default(''),
  type: z.enum(['income', 'expense'], { message: 'type must be income or expense' }),
  amount: z.number().min(1, 'amount must be at least 1'),
  description: z.string().optional().default(''),
  date: z.string().min(1, 'date is required'),
});

export const overtimeSchema = z.object({
  date: z.string().min(1, 'date is required'),
  hours: z.number().min(0.5, 'hours must be at least 0.5').max(12, 'hours must not exceed 12'),
  is_holiday: z.boolean().optional().default(false),
  notes: z.string().optional().default(''),
});

export const disburseSchema = z.object({
  period_start: z.string().min(1, 'period_start is required'),
  period_end: z.string().min(1, 'period_end is required'),
  wallet_id: z.string().uuid('wallet_id is required'),
});

export const goalSchema = z.object({
  name: z.string().min(1, 'name is required'),
  target_amount: z.number().min(1, 'target_amount must be at least 1'),
  deadline: z.string().nullable().optional(),
  tracking_mode: z.enum(['manual', 'single_wallet', 'multiple_wallet', 'all_wallet']).optional().default('manual'),
  notes: z.string().optional().default(''),
  wallet_ids: z.array(z.string()).optional().default([]),
});

export const incomeSchema = z.object({
  wallet_id: z.string().optional().default(''),
  amount: z.number().min(1, 'amount must be at least 1'),
  source: z.string().min(1, 'source is required'),
  description: z.string().optional().default(''),
  date: z.string().min(1, 'date is required'),
});

export const dailyBudgetSchema = z.object({
  mode: z.enum(['manual', 'formula']),
  manual_amount: z.number().optional().default(0),
  formula_wallet_id: z.string().optional().default(''),
  formula_days_remaining: z.number().optional().default(30),
});

export const profileSchema = z.object({
  name: z.string().min(1, 'name is required'),
  phone: z.string().optional().default(''),
  telegram: z.string().optional().default(''),
  salary: z.number().optional().default(0),
  meal_allowance: z.number().optional().default(0),
});

export const changePasswordSchema = z.object({
  old_password: z.string().min(1, 'old_password is required'),
  new_password: z.string().min(6, 'new_password must be at least 6 characters'),
});

export const adminCreateUserSchema = z.object({
  name: z.string().min(1, 'name is required'),
  email: z.string().email('email must be a valid email address'),
  password: z.string().min(6, 'password must be at least 6 characters'),
  role: z.string().optional().default('user'),
  salary: z.number().optional().default(0),
  meal_allowance: z.number().optional().default(0),
});

export const resetPasswordSchema = z.object({
  password: z.string().min(6, 'password must be at least 6 characters'),
});
