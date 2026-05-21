import { registerSchema, loginSchema, walletSchema, transactionSchema, overtimeSchema, goalSchema } from '../src/utils/validators';

describe('validators', () => {
  describe('registerSchema', () => {
    it('should pass with valid input', () => {
      const result = registerSchema.safeParse({ name: 'Test', email: 'test@example.com', password: 'password123' });
      expect(result.success).toBe(true);
    });

    it('should fail with empty name', () => {
      const result = registerSchema.safeParse({ name: '', email: 'test@example.com', password: 'password123' });
      expect(result.success).toBe(false);
    });

    it('should fail with invalid email', () => {
      const result = registerSchema.safeParse({ name: 'Test', email: 'notanemail', password: 'password123' });
      expect(result.success).toBe(false);
    });

    it('should fail with short password', () => {
      const result = registerSchema.safeParse({ name: 'Test', email: 'test@example.com', password: '123' });
      expect(result.success).toBe(false);
    });
  });

  describe('loginSchema', () => {
    it('should pass with valid input', () => {
      const result = loginSchema.safeParse({ email: 'test@example.com', password: 'password' });
      expect(result.success).toBe(true);
    });

    it('should fail with empty email', () => {
      const result = loginSchema.safeParse({ email: '', password: 'password' });
      expect(result.success).toBe(false);
    });
  });

  describe('walletSchema', () => {
    it('should pass with valid input', () => {
      const result = walletSchema.safeParse({ name: 'My Wallet' });
      expect(result.success).toBe(true);
    });

    it('should fail with empty name', () => {
      const result = walletSchema.safeParse({ name: '' });
      expect(result.success).toBe(false);
    });
  });

  describe('transactionSchema', () => {
    it('should pass with valid input', () => {
      const result = transactionSchema.safeParse({
        wallet_id: '550e8400-e29b-41d4-a716-446655440000',
        type: 'income',
        amount: 50000,
        date: '2024-01-01',
      });
      expect(result.success).toBe(true);
    });

    it('should fail with invalid type', () => {
      const result = transactionSchema.safeParse({
        wallet_id: '550e8400-e29b-41d4-a716-446655440000',
        type: 'other',
        amount: 50000,
        date: '2024-01-01',
      });
      expect(result.success).toBe(false);
    });

    it('should fail with zero amount', () => {
      const result = transactionSchema.safeParse({
        wallet_id: '550e8400-e29b-41d4-a716-446655440000',
        type: 'expense',
        amount: 0,
        date: '2024-01-01',
      });
      expect(result.success).toBe(false);
    });
  });

  describe('overtimeSchema', () => {
    it('should pass with valid input', () => {
      const result = overtimeSchema.safeParse({ date: '2024-01-01', hours: 2 });
      expect(result.success).toBe(true);
    });

    it('should fail with hours below 0.5', () => {
      const result = overtimeSchema.safeParse({ date: '2024-01-01', hours: 0.2 });
      expect(result.success).toBe(false);
    });

    it('should fail with hours above 12', () => {
      const result = overtimeSchema.safeParse({ date: '2024-01-01', hours: 13 });
      expect(result.success).toBe(false);
    });
  });

  describe('goalSchema', () => {
    it('should pass with valid input', () => {
      const result = goalSchema.safeParse({ name: 'Buy Car', target_amount: 200000000 });
      expect(result.success).toBe(true);
    });

    it('should fail with zero target', () => {
      const result = goalSchema.safeParse({ name: 'Buy Car', target_amount: 0 });
      expect(result.success).toBe(false);
    });
  });
});
