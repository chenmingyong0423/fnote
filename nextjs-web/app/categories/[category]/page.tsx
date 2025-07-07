"use client";

import ArticleList from "@/src/components/ArticleList";
import { notFound, useSearchParams } from "next/navigation";
import React from "react";

export default function CategoryPage({ params }: { params: Promise<{ category: string }> }) {
  const searchParams = useSearchParams();
  const field = (searchParams?.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const pageSize = Number(searchParams?.get("pageSize") || 10);
  const { category } = React.use(params);
  if (!category) return notFound();
  return <ArticleList category={category} field={field} page={1} pageSize={pageSize} />;
}
