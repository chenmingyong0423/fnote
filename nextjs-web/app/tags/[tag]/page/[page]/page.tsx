"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { useSearchParams } from "next/navigation";
import { notFound } from "next/navigation";

export default function TagPageWithPagination({ params }: { params: { tag: string; page: string } }) {
  const searchParams = useSearchParams();
  const field = (searchParams.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const page = Number(params.page || 1);
  const pageSize = Number(searchParams.get("pageSize") || 10);
  if (!params.tag) return notFound();
  return <ArticleListContainer tag={params.tag} field={field} page={page} pageSize={pageSize} />;
}
