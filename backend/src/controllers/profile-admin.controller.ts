import { Response } from 'express';
import { profileService, adminService } from '../services';
import { AuthRequest } from '../middleware/auth';
import * as response from '../utils/response';
import { profileSchema, changePasswordSchema, adminCreateUserSchema, resetPasswordSchema } from '../utils/validators';
import { parsePagination, totalPages } from '../utils/pagination';
import { ZodError } from 'zod';

function formatZodErrors(err: ZodError): Record<string, string> {
  const errors: Record<string, string> = {};
  for (const issue of err.issues) {
    errors[issue.path.join('.')] = issue.message;
  }
  return errors;
}

export const profileController = {
  async update(req: AuthRequest, res: Response) {
    try {
      const input = profileSchema.parse(req.body);
      const user = await profileService.update(req.userId!, input);
      response.ok(res, 'Profile updated', user);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async changePassword(req: AuthRequest, res: Response) {
    try {
      const input = changePasswordSchema.parse(req.body);
      await profileService.changePassword(req.userId!, input.old_password, input.new_password);
      response.ok(res, 'Password changed successfully');
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },
};

export const adminController = {
  async listUsers(req: AuthRequest, res: Response) {
    try {
      const { perPage, offset, page } = parsePagination(req);
      const { users, total } = await adminService.listUsers(perPage, offset);
      response.successWithMeta(res, 200, 'Users', users, {
        page, per_page: perPage, total, total_pages: totalPages(total, perPage),
      });
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async createUser(req: AuthRequest, res: Response) {
    try {
      const input = adminCreateUserSchema.parse(req.body);
      const user = await adminService.createUser(input);
      response.created(res, 'User created', user);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 409, err.message);
    }
  },

  async updateUser(req: AuthRequest, res: Response) {
    try {
      const input = adminCreateUserSchema.parse(req.body);
      const user = await adminService.updateUser(req.params.id, input);
      response.ok(res, 'User updated', user);
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async deleteUser(req: AuthRequest, res: Response) {
    try {
      await adminService.deleteUser(req.params.id);
      response.ok(res, 'User deleted');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },

  async resetPassword(req: AuthRequest, res: Response) {
    try {
      const input = resetPasswordSchema.parse(req.body);
      await adminService.resetPassword(req.params.id, input.password);
      response.ok(res, 'Password reset successful');
    } catch (err: any) {
      if (err instanceof ZodError) return response.validationError(res, formatZodErrors(err));
      response.error(res, 400, err.message);
    }
  },

  async getSettings(_req: AuthRequest, res: Response) {
    try {
      const settings = await adminService.getSettings();
      response.ok(res, 'Site settings', settings);
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },

  async updateSettings(req: AuthRequest, res: Response) {
    try {
      await adminService.updateSettings(req.body);
      response.ok(res, 'Settings updated');
    } catch (err: any) {
      response.error(res, 400, err.message);
    }
  },

  async getActivityLogs(req: AuthRequest, res: Response) {
    try {
      const { perPage, offset, page } = parsePagination(req, 50);
      const { logs, total } = await adminService.getActivityLogs(perPage, offset);
      response.successWithMeta(res, 200, 'Activity logs', logs, {
        page, per_page: perPage, total, total_pages: totalPages(total, perPage),
      });
    } catch (err: any) {
      response.internalError(res, err.message);
    }
  },
};
