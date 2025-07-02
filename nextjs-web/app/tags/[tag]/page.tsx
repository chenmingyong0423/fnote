"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound, useSearchParams } from "next/navigation";
import React from "react";

export default function TagPage({ params }: { params: Promise<{ tag: string }> }) {
  const searchParams = useSearchParams();
  const field = (searchParams?.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const pageSize = Number(searchParams?.get("pageSize") || 10);
  const { tag } = React.use(params);
  if (!tag) return notFound();
  return <ArticleListContainer tag={tag} field={field} page={1} pageSize={pageSize} />;
}
