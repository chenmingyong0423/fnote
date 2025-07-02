import PostDetail from "@/src/components/PostDetail";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";

export default async function AboutPage() {
  const res = await getPostDetail("about-me");
  if (res.code !== 0 || !res.data) return notFound();
  return <PostDetail post={res.data} />;
}