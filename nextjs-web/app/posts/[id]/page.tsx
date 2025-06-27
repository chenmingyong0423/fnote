import PostDetail from "@/src/components/PostDetail";

export default function PostDetailPage({ params }: { params: { id: string } }) {
  return <PostDetail id={params.id} />;
}
