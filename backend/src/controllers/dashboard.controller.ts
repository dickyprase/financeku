import { Response } from 'express';
import { dashboardService, dailyBudgetService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { dailyBudgetSchema } from '../utils/validators';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const dashboardController = {
  async getSummary(req: AuthRequest, res: Response) {
    try {
      const summary = await dashboardService.getSummary(req.userId!);
      response.ok(res, 'Dashboard summary', summary);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async getCashflowReport(req: AuthRequest, res: Response) {
    try {
      const month = req.query.month as string;
      if (!month) {
        return response.error(res, 400, 'month parameter is required (format: YYYY-MM)');
      }
      const report = await dashboardService.getCashflowReport(req.userId!, month);
      response.ok(res, 'Cashflow report', report);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },
};

export const dailyBudgetController = {
  async get(req: AuthRequest, res: Response) {
    try {
      const setting = await dailyBudgetService.get(req.userId!);
      response.ok(res, 'Daily budget settings', setting);
    } catch (err: any) {
      response.ok(res, 'Daily budget settings', null);
    }
  },

  async update(req: AuthRequest, res: Response) {
    try {
      const input = dailyBudgetSchema.parse(req.body);
      const setting = await dailyBudgetService.upsert(req.userId!, input);
      response.ok(res, 'Daily budget settings updated', setting);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.internalError(res, err.message);
    }
  },

  async getToday(req: AuthRequest, res: Response) {
    try {
      const today = await dailyBudgetService.getToday(req.userId!);
      response.ok(res, "Today's budget", today);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },
};
