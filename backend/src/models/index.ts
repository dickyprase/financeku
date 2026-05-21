export interface User {
  id: string;
  name: string;
  email: string;
  password: string;
  role: string;
  phone: string | null;
  telegram: string | null;
  salary: number;
  meal_allowance: number;
  is_active: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface Wallet {
  id: string;
  user_id: string;
  name: string;
  balance: number;
  icon: string | null;
  color: string | null;
  is_active: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface Category {
  id: string;
  user_id: string | null;
  name: string;
  type: 'income' | 'expense';
  icon: string | null;
  color: string | null;
  budget_limit: number;
  is_default: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface Transaction {
  id: string;
  user_id: string;
  wallet_id: string;
  category_id: string | null;
  type: 'income' | 'expense' | 'transfer';
  amount: number;
  description: string | null;
  date: string;
  reference_id: string | null;
  reference_type: string | null;
  created_at: Date;
  updated_at: Date;
}

export interface OvertimeRecord {
  id: string;
  user_id: string;
  date: string;
  hours: number;
  is_holiday: boolean;
  amount: number;
  meal_amount: number;
  total_amount: number;
  period_start: string | null;
  period_end: string | null;
  is_disbursed: boolean;
  disbursed_at: Date | null;
  notes: string | null;
  created_at: Date;
  updated_at: Date;
}

export interface Income {
  id: string;
  user_id: string;
  wallet_id: string | null;
  amount: number;
  source: string;
  description: string | null;
  date: string;
  is_from_overtime: boolean;
  overtime_period_start: string | null;
  overtime_period_end: string | null;
  created_at: Date;
  updated_at: Date;
}

export interface Goal {
  id: string;
  user_id: string;
  name: string;
  target_amount: number;
  current_amount: number;
  deadline: string | null;
  status: 'active' | 'completed' | 'archived';
  tracking_mode: 'manual' | 'single_wallet' | 'multiple_wallet' | 'all_wallet';
  notes: string | null;
  created_at: Date;
  updated_at: Date;
}

export interface GoalWallet {
  id: string;
  goal_id: string;
  wallet_id: string;
  created_at: Date;
}

export interface DailyBudgetSetting {
  id: string;
  user_id: string;
  mode: 'manual' | 'formula';
  manual_amount: number;
  formula_wallet_id: string | null;
  formula_days_remaining: number;
  exclude_categories: string[] | null;
  created_at: Date;
  updated_at: Date;
}

export interface ActivityLog {
  id: string;
  user_id: string | null;
  action: string;
  entity_type: string | null;
  entity_id: string | null;
  details: any;
  ip_address: string | null;
  created_at: Date;
}

export interface SiteSetting {
  id: string;
  key: string;
  value: string | null;
  description: string | null;
  updated_at: Date;
}
