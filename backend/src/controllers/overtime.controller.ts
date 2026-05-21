import { Response } from 'express';
import { overtimeService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { overtimeSchema, disburseSchema } from '../utils/validators';
import { parsePagination, totalPages } from '../utils/pagination';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const overtimeController = {
  async create(req: AuthRequest, res: Response) {
    try {
      const input = overtimeSchema.parse(req.body);
      const record = await overtimeService.create(req.userId!, input);
      response.created(res, 'Overtime record created', record);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 500, err.message);
    }
  },

  async list(req: AuthRequest, res: Response) {
    try {
      const { page, perPage, offset } = parsePagination(req);
      const month = (req.query.month as string) || '';
      const { records, total } = await overtimeService.list(req.userId!, month, perPage, offset);
      response.successWithMeta(res, 200, 'Overtime records', records, {
        page, per_page: perPage, total, total_pages: totalPages(total, perPage),
      });
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async getById(req: AuthRequest, res: Response) {
    try {
      const record = await overtimeService.getById(req.userId!, req.params.id);
      if (!record) return response.notFound(res, 'Overtime record not found');
      response.ok(res, 'Overtime record', record);
    } catch (err: any) {
      response.notFound(res, err.message);
    }
  },

  async update(req: AuthRequest, res: Response) {
    try {
      const input = overtimeSchema.parse(req.body);
      const record = await overtimeService.update(req.userId!, req.params.id, input);
      response.ok(res, 'Overtime record updated', record);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async delete(req: AuthRequest, res: Response) {
    try {
      await overtimeService.delete(req.userId!, req.params.id);
      response.ok(res, 'Overtime record deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },

  async calculate(req: AuthRequest, res: Response) {
    try {
      const hours = parseFloat(req.query.hours as string) || 0;
      const isHoliday = req.query.is_holiday === 'true' || req.query.is_holiday === '1';
      const calc = await overtimeService.calculate(req.userId!, hours, isHoliday);
      response.ok(res, 'Overtime calculation', calc);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async disburse(req: AuthRequest, res: Response) {
    try {
      const input = disburseSchema.parse(req.body);
      await overtimeService.disburse(req.userId!, input);
      response.ok(res, 'Period disbursed successfully');
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },
};
