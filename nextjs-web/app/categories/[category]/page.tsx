"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound, useSearchParams } from "next/navigation";
import React from "react";

export default function CategoryPage({ params }: { params: Promise<{ category: string }> }) {
  const searchParams = useSearchParams();
  const field = (searchParams.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const pageSize = Number(searchParams.get("pageSize") || 10);
  const { category } = React.use(params);
  if (!category) return notFound();
  return <ArticleListContainer category={category} field={field} page={1} pageSize={pageSize} />;
}
