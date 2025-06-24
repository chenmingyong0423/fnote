"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { useSearchParams } from "next/navigation";
import { notFound } from "next/navigation";

export default function CategoryPageWithPagination({ params }: { params: { category: string; page: string } }) {
  const searchParams = useSearchParams();
  const field = (searchParams.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const page = Number(params.page || 1);
  const pageSize = Number(searchParams.get("pageSize") || 10);
  if (!params.category) return notFound();
  return <ArticleListContainer category={params.category} field={field} page={page} pageSize={pageSize} />;
}
