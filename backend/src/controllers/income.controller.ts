import { Response } from 'express';
import { incomeService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { incomeSchema } from '../utils/validators';
import { parsePagination, totalPages } from '../utils/pagination';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const incomeController = {
  async create(req: AuthRequest, res: Response) {
    try {
      const input = incomeSchema.parse(req.body);
      const income = await incomeService.create(req.userId!, input);
      response.created(res, 'Income created', income);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async list(req: AuthRequest, res: Response) {
    try {
      const { page, perPage, offset } = parsePagination(req);
      const { incomes, total } = await incomeService.list(req.userId!, perPage, offset);
      response.successWithMeta(res, 200, 'Incomes', incomes, {
        page, per_page: perPage, total, total_pages: totalPages(total, perPage),
      });
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async delete(req: AuthRequest, res: Response) {
    try {
      await incomeService.delete(req.userId!, req.params.id);
      response.ok(res, 'Income deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },
};
