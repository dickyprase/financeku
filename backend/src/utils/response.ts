import { Response } from 'express';

export interface Meta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export function success(res: Response, status: number, message: string, data?: any) {
  res.status(status).json({ success: true, message, data });
}

export function successWithMeta(res: Response, status: number, message: string, data: any, meta: Meta) {
  res.status(status).json({ success: true, message, data, meta });
}

export function error(res: Response, status: number, message: string) {
  res.status(status).json({ success: false, message });
}

export function validationError(res: Response, errors: Record<string, string>) {
  res.status(422).json({ success: false, message: 'Validation failed', errors });
}

export function created(res: Response, message: string, data?: any) {
  success(res, 201, message, data);
}

export function ok(res: Response, message: string, data?: any) {
  success(res, 200, message, data);
}

export function unauthorized(res: Response, message = 'Unauthorized') {
  error(res, 401, message);
}

export function forbidden(res: Response, message = 'Forbidden') {
  error(res, 403, message);
}

export function notFound(res: Response, message = 'Resource not found') {
  error(res, 404, message);
}

export function internalError(res: Response, message = 'Internal server error') {
  error(res, 500, message);
}
