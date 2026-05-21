import { Response } from 'express';
import { walletService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { walletSchema, transferSchema } from '../utils/validators';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const walletController = {
  async create(req: AuthRequest, res: Response) {
    try {
      const input = walletSchema.parse(req.body);
      const wallet = await walletService.create(req.userId!, input);
      response.created(res, 'Wallet created', wallet);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.internalError(res, err.message);
    }
  },

  async list(req: AuthRequest, res: Response) {
    try {
      const wallets = await walletService.list(req.userId!);
      response.ok(res, 'Wallets', wallets);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async getById(req: AuthRequest, res: Response) {
    try {
      const wallet = await walletService.getById(req.userId!, req.params.id);
      if (!wallet) return response.notFound(res, 'Wallet not found');
      response.ok(res, 'Wallet', wallet);
    } catch (err: any) {
      response.notFound(res, err.message);
    }
  },

  async update(req: AuthRequest, res: Response) {
    try {
      const input = walletSchema.parse(req.body);
      const wallet = await walletService.update(req.userId!, req.params.id, input);
      response.ok(res, 'Wallet updated', wallet);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async delete(req: AuthRequest, res: Response) {
    try {
      await walletService.delete(req.userId!, req.params.id);
      response.ok(res, 'Wallet deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },

  async transfer(req: AuthRequest, res: Response) {
    try {
      const input = transferSchema.parse(req.body);
      await walletService.transfer(req.userId!, input);
      response.ok(res, 'Transfer successful');
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },
};
