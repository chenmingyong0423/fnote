import ArticleListContainer from "@/src/components/ArticleListContainer";

export default function SearchPage({ searchParams }: { searchParams: { q?: string; field?: "latest" | "oldest" | "likes" } }) {
  return <ArticleListContainer keyword={searchParams.q} field={searchParams.field || "latest"} />;
}
