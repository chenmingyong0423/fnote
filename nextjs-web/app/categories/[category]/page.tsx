import ArticleListContainer from "@/src/components/ArticleListContainer";
import { notFound } from "next/navigation";

export default function CategoryPage({ params }: { params: { category: string } }) {
  if (!params.category) return notFound();
  return <ArticleListContainer category={params.category} field="latest" />;
}
