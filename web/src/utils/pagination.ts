import { notFound } from "next/navigation";

export type ListSearchParams = Record<string, string | string[] | undefined>;

function single(value: string | string[] | undefined) {
  if (Array.isArray(value)) notFound();
  return value;
}

export function resolvePagination(base: string, rawPage: string | undefined, query: ListSearchParams) {
  // 只支持列表根路径和 /page/2 等标准路径，不兼容查询参数分页。
  if (query.page !== undefined) notFound();
  const value = rawPage ?? "1";
  if (!/^\d+$/.test(value)) notFound();
  const page = Number(value);
  if (!Number.isSafeInteger(page) || page < 1) notFound();
  if (rawPage !== undefined && (page === 1 || rawPage !== String(page))) notFound();
  const sizeValue = single(query.pageSize);
  if (sizeValue !== undefined && !/^\d+$/.test(sizeValue)) notFound();
  const pageSize = sizeValue === undefined ? 10 : Number(sizeValue);
  if (![5, 10, 20, 50].includes(pageSize)) notFound();
  const filter = single(query.filter) ?? "latest";
  if (!["latest", "oldest", "likes"].includes(filter)) notFound();
  const field = filter as "latest" | "oldest" | "likes";
  const pathname = page === 1 ? base : `${base}/page/${page}`;

  return { page, pageSize, field, pathname };
}

export function ensurePageExists(page: number, totalPages: number) {
  // 空列表的第一页仍可展示；其余不存在的分页必须返回 404。
  if (page > 1 && page > totalPages) notFound();
}
