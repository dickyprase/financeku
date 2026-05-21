import { Request } from 'express';

export interface PaginationParams {
  page: number;
  perPage: number;
  offset: number;
}

export function parsePagination(req: Request, defaultPerPage = 20, maxPerPage = 100): PaginationParams {
  let page = parseInt(req.query.page as string) || 1;
  let perPage = parseInt(req.query.per_page as string) || defaultPerPage;

  if (page < 1) page = 1;
  if (perPage < 1) perPage = defaultPerPage;
  if (perPage > maxPerPage) perPage = maxPerPage;

  return {
    page,
    perPage,
    offset: (page - 1) * perPage,
  };
}

export function totalPages(total: number, perPage: number): number {
  return Math.ceil(total / perPage);
}
