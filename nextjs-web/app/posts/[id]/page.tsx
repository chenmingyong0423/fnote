import PostDetail from "@/src/components/PostDetail";

export default async function PostDetailPage({ params }: { params: { id: string } }) {
  const {id} = await params;
  return <PostDetail id={id} />;
}
