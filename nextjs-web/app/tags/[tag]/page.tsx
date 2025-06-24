"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound, useSearchParams } from "next/navigation";

export default function TagPage({ params }: { params: { tag: string } }) {
  const searchParams = useSearchParams();
  const field = (searchParams.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const pageSize = Number(searchParams.get("pageSize") || 10);
  if (!params.tag) return notFound();
  return <ArticleListContainer tag={params.tag} field={field} page={1} pageSize={pageSize} />;
}
