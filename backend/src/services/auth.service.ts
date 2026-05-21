import jwt from 'jsonwebtoken';
import { config } from '../config';
import { userRepository } from '../repositories';
import { hashPassword, checkPassword } from '../utils/hash';

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export const authService = {
  async register(name: string, email: string, password: string) {
    const existing = await userRepository.findByEmail(email);
    if (existing) {
      throw new Error('email already registered');
    }

    const hashed = await hashPassword(password);
    const user = await userRepository.create({
      name,
      email,
      password: hashed,
      role: 'user',
    });

    return user;
  },

  async login(email: string, password: string) {
    const user = await userRepository.findByEmail(email);
    if (!user) {
      throw new Error('invalid email or password');
    }

    if (!user.is_active) {
      throw new Error('account is deactivated');
    }

    const valid = await checkPassword(password, user.password);
    if (!valid) {
      throw new Error('invalid email or password');
    }

    const tokens = generateTokenPair(user.id, user.email, user.role);
    return { user, tokens };
  },

  async refreshToken(refreshTokenStr: string): Promise<TokenPair> {
    try {
      const payload = jwt.verify(refreshTokenStr, config.jwt.secret) as any;

      if (payload.type !== 'refresh') {
        throw new Error('invalid token type');
      }

      const user = await userRepository.findById(payload.user_id);
      if (!user) {
        throw new Error('user not found');
      }

      if (!user.is_active) {
        throw new Error('account is deactivated');
      }

      return generateTokenPair(user.id, user.email, user.role);
    } catch {
      throw new Error('invalid refresh token');
    }
  },

  async getCurrentUser(userId: string) {
    const user = await userRepository.findById(userId);
    if (!user) {
      throw new Error('user not found');
    }
    return user;
  },
};

function generateTokenPair(userId: string, email: string, role: string): TokenPair {
  const accessExp = config.jwt.accessExpMinutes * 60;
  const refreshExp = config.jwt.refreshExpDays * 24 * 60 * 60;

  const accessToken = jwt.sign(
    { user_id: userId, email, role, type: 'access' },
    config.jwt.secret,
    { expiresIn: accessExp }
  );

  const refreshToken = jwt.sign(
    { user_id: userId, type: 'refresh' },
    config.jwt.secret,
    { expiresIn: refreshExp }
  );

  return {
    access_token: accessToken,
    refresh_token: refreshToken,
    expires_in: accessExp,
  };
}
