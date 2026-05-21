import { Response } from 'express';
import { transactionService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { transactionSchema } from '../utils/validators';
import { parsePagination, totalPages } from '../utils/pagination';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const transactionController = {
  async create(req: AuthRequest, res: Response) {
    try {
      const input = transactionSchema.parse(req.body);
      const tx = await transactionService.create(req.userId!, input);
      response.created(res, 'Transaction created', tx);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async list(req: AuthRequest, res: Response) {
    try {
      const { page, perPage, offset } = parsePagination(req);
      const filter = {
        walletId: (req.query.wallet_id as string) || undefined,
        categoryId: (req.query.category_id as string) || undefined,
        type: (req.query.type as string) || undefined,
        dateFrom: (req.query.date_from as string) || undefined,
        dateTo: (req.query.date_to as string) || undefined,
        limit: perPage,
        offset,
      };
      const { transactions, total } = await transactionService.list(req.userId!, filter);
      response.successWithMeta(res, 200, 'Transactions', transactions, {
        page, per_page: perPage, total, total_pages: totalPages(total, perPage),
      });
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async getById(req: AuthRequest, res: Response) {
    try {
      const tx = await transactionService.getById(req.userId!, req.params.id);
      if (!tx) return response.notFound(res, 'Transaction not found');
      response.ok(res, 'Transaction', tx);
    } catch (err: any) {
      response.notFound(res, err.message);
    }
  },

  async update(req: AuthRequest, res: Response) {
    try {
      const input = transactionSchema.parse(req.body);
      const tx = await transactionService.update(req.userId!, req.params.id, input);
      response.ok(res, 'Transaction updated', tx);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async delete(req: AuthRequest, res: Response) {
    try {
      await transactionService.delete(req.userId!, req.params.id);
      response.ok(res, 'Transaction deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },
};
