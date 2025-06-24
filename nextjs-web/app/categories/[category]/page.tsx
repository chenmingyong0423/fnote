"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound, useSearchParams } from "next/navigation";

export default function CategoryPage({ params }: { params: { category: string } }) {
  const searchParams = useSearchParams();
  const field = (searchParams.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const page = Number(searchParams.get("page") || 1);
  const pageSize = Number(searchParams.get("pageSize") || 10);
  if (!params.category) return notFound();
  return <ArticleListContainer category={params.category} field={field} page={page} pageSize={pageSize} />;
}
