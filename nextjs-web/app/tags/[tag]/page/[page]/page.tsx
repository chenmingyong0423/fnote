"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { useSearchParams } from "next/navigation";
import { notFound } from "next/navigation";
import React from "react";

export default function TagPageWithPagination({ params }: { params: Promise<{ tag: string; page: string }> }) {
  const searchParams = useSearchParams();
  const { tag, page } = React.use(params);
  const field = (searchParams?.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const pageNumber = Number(page || 1);
  const pageSize = Number(searchParams?.get("pageSize") || 10);
  if (!tag) return notFound();
  return <ArticleListContainer tag={tag} field={field} page={pageNumber} pageSize={pageSize} />;
}
