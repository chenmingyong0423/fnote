"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound } from "next/navigation";
import { useSearchParams } from "next/navigation";

export default function TagPage({ params }: { params: { tag: string } }) {
  const searchParams = useSearchParams();
  const field = (searchParams.get("field") as "latest" | "oldest" | "likes") || "latest";
  if (!params.tag) return notFound();
  return <ArticleListContainer tag={params.tag} field={field} />;
}
