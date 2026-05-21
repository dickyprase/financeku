import { hashPassword, checkPassword } from '../src/utils/hash';

describe('hash', () => {
  it('should hash a password', async () => {
    const password = 'testpassword123';
    const hashed = await hashPassword(password);

    expect(hashed).not.toBe('');
    expect(hashed).not.toBe(password);
  });

  it('should verify correct password', async () => {
    const password = 'testpassword123';
    const hashed = await hashPassword(password);

    const result = await checkPassword(password, hashed);
    expect(result).toBe(true);
  });

  it('should reject wrong password', async () => {
    const hashed = await hashPassword('testpassword123');

    const result = await checkPassword('wrongpassword', hashed);
    expect(result).toBe(false);
  });

  it('should reject empty password', async () => {
    const hashed = await hashPassword('somepassword');

    const result = await checkPassword('', hashed);
    expect(result).toBe(false);
  });
});
