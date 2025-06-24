import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound } from "next/navigation";

export default function TagPage({ params }: { params: { tag: string } }) {
  if (!params.tag) return notFound();
  return <ArticleListContainer tag={params.tag} field="latest" />;
}
