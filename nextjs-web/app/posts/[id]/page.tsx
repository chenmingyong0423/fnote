import PostDetail from "@/src/components/PostDetail";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";

type Params = Promise<{
    id: string;
  }>

export default async function PostDetailPage({ params }: { params: Params }) {
  const { id } = await params
  const res = await getPostDetail(id);
  if (res.code !== 0 || !res.data) return notFound();
  return <PostDetail post={res.data} />;
}
