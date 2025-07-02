"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { useSearchParams } from "next/navigation";
import { notFound } from "next/navigation";
import React from "react";

export default function CategoryPageWithPagination({ params }: { params: Promise<{ category: string; page: string }> }) {
  const searchParams = useSearchParams();
  const { category, page } = React.use(params);
  const field = (searchParams?.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const pageNumber = Number(page || 1);
  const pageSize = Number(searchParams?.get("pageSize") || 10);
  if (!category) return notFound();
  return <ArticleListContainer category={category} field={field} page={pageNumber} pageSize={pageSize} />;
}
