import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';
import { config } from '../config';
import * as response from '../utils/response';

export interface AuthRequest extends Request {
  userId?: string;
  userRole?: string;
}

export function requireAuth(req: AuthRequest, res: Response, next: NextFunction) {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    return response.unauthorized(res, 'Missing authorization header');
  }

  const parts = authHeader.split(' ');
  if (parts.length !== 2 || parts[0] !== 'Bearer') {
    return response.unauthorized(res, 'Invalid authorization header format');
  }

  try {
    const payload = jwt.verify(parts[1], config.jwt.secret) as any;

    if (payload.type !== 'access') {
      return response.unauthorized(res, 'Invalid token type');
    }

    req.userId = payload.user_id;
    req.userRole = payload.role;
    next();
  } catch {
    return response.unauthorized(res, 'Invalid or expired token');
  }
}

export function requireAdmin(req: AuthRequest, res: Response, next: NextFunction) {
  if (req.userRole !== 'admin') {
    return response.forbidden(res, 'Admin access required');
  }
  next();
}
