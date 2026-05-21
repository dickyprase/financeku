import { Response } from 'express';
import { categoryService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { categorySchema } from '../utils/validators';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const categoryController = {
  async create(req: AuthRequest, res: Response) {
    try {
      const input = categorySchema.parse(req.body);
      const cat = await categoryService.create(req.userId!, input);
      response.created(res, 'Category created', cat);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.internalError(res, err.message);
    }
  },

  async list(req: AuthRequest, res: Response) {
    try {
      const type = (req.query.type as string) || undefined;
      const categories = await categoryService.list(req.userId!, type);
      response.ok(res, 'Categories', categories);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async update(req: AuthRequest, res: Response) {
    try {
      const input = categorySchema.parse(req.body);
      const cat = await categoryService.update(req.userId!, req.params.id, input);
      response.ok(res, 'Category updated', cat);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async delete(req: AuthRequest, res: Response) {
    try {
      await categoryService.delete(req.userId!, req.params.id);
      response.ok(res, 'Category deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },
};
