import { Response } from 'express';
import { goalService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { goalSchema } from '../utils/validators';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const goalController = {
  async create(req: AuthRequest, res: Response) {
    try {
      const input = goalSchema.parse(req.body);
      const goal = await goalService.create(req.userId!, input);
      response.created(res, 'Goal created', goal);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.internalError(res, err.message);
    }
  },

  async list(req: AuthRequest, res: Response) {
    try {
      const goals = await goalService.list(req.userId!);
      response.ok(res, 'Goals', goals);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async getById(req: AuthRequest, res: Response) {
    try {
      const goal = await goalService.getById(req.userId!, req.params.id);
      if (!goal) return response.notFound(res, 'Goal not found');
      response.ok(res, 'Goal', goal);
    } catch (err: any) {
      response.notFound(res, err.message);
    }
  },

  async update(req: AuthRequest, res: Response) {
    try {
      const input = goalSchema.parse(req.body);
      const goal = await goalService.update(req.userId!, req.params.id, input);
      response.ok(res, 'Goal updated', goal);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async delete(req: AuthRequest, res: Response) {
    try {
      await goalService.delete(req.userId!, req.params.id);
      response.ok(res, 'Goal deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },

  async getProgress(req: AuthRequest, res: Response) {
    try {
      const progress = await goalService.getProgress(req.userId!, req.params.id);
      response.ok(res, 'Goal progress', progress);
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },
};
