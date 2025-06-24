"use client";

import ArticleListContainer from "@/src/components/ArticleListContainer";
import { useSearchParams } from "next/navigation";

export default function SearchPage() {
  const params = useSearchParams();
  const field = (params.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const page = Number(params.get("page") || 1);
  const pageSize = Number(params.get("pageSize") || 10);
  const q = params.get("q") || undefined;
  return <ArticleListContainer keyword={q} field={field} page={page} pageSize={pageSize} />;
}
