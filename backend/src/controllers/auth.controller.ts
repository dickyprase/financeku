import { Request, Response } from 'express';
import { authService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { registerSchema, loginSchema, refreshSchema } from '../utils/validators';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    const field = issue.path.join('.');
    errors[field] = issue.message;
  }
  return errors;
}

export const authController = {
  async register(req: Request, res: Response) {
    try {
      const input = registerSchema.parse(req.body);
      const user = await authService.register(input.name, input.email, input.password);
      response.created(res, 'Registration successful', user);
    } catch (err: any) {
      if (err instanceof ZodError) {
        return response.validationError(res, formatZodErrors(err));
      }
      response.error(res, 409, err.message);
    }
  },

  async login(req: Request, res: Response) {
    try {
      const input = loginSchema.parse(req.body);
      const { user, tokens } = await authService.login(input.email, input.password);
      response.ok(res, 'Login successful', { user, tokens });
    } catch (err: any) {
      if (err instanceof ZodError) {
        return response.validationError(res, formatZodErrors(err));
      }
      response.error(res, 401, err.message);
    }
  },

  async refresh(req: Request, res: Response) {
    try {
      const input = refreshSchema.parse(req.body);
      const tokens = await authService.refreshToken(input.refresh_token);
      response.ok(res, 'Token refreshed', tokens);
    } catch (err: any) {
      if (err instanceof ZodError) {
        return response.validationError(res, formatZodErrors(err));
      }
      response.unauthorized(res, err.message);
    }
  },

  async logout(_req: Request, res: Response) {
    response.ok(res, 'Logged out successfully');
  },

  async me(req: AuthRequest, res: Response) {
    try {
      const user = await authService.getCurrentUser(req.userId!);
      response.ok(res, 'User profile', user);
    } catch (err: any) {
      response.notFound(res, err.message);
    }
  },
};
