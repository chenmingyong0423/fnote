import PostDetail from "@/src/components/PostDetail";
import { getPostDetailOrNull } from "@/src/api/posts";
import { notFound } from "next/navigation";
import { getCommonConfig } from "@/src/api/config";
import type { Metadata } from "next";
import { getCommentsByPostId } from "@/src/api/comments";
import { buildBlogPostingJsonLd, buildPageMetadata, getPostPath, serializeJsonLd } from "@/src/utils/seo";

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }): Promise<Metadata> {
  const { id } = await params;
  const post = await getPostDetailOrNull(id);
  if (!post) notFound();
  const config = await getCommonConfig();
  const metadata = buildPageMetadata(config, {
    title: post.title,
    description: post.meta_description || post.summary,
    pathname: getPostPath(post._id),
    image: post.cover_img,
  });
  return {
    ...metadata,
    keywords: post.meta_keywords,
    authors: [{ name: post.author }],
    openGraph: {
      ...metadata.openGraph,
      type: "article",
      publishedTime: new Date(post.created_at * 1000).toISOString(),
      modifiedTime: new Date((post.updated_at || post.created_at) * 1000).toISOString(),
      authors: [post.author],
      tags: post.tags.map((tag) => tag.name),
    },
  };
}

type Params = Promise<{
    id: string;
  }>

export default async function PostDetailPage({ params }: { params: Params }) {
  const { id } = await params
  const post = await getPostDetailOrNull(id);
  if (!post) return notFound();

  const [comments, config] = await Promise.all([
    getCommentsByPostId(post._id).catch(() => []),
    getCommonConfig(),
  ]);
  const blogPostingJsonLd = buildBlogPostingJsonLd(post, config);

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: serializeJsonLd(blogPostingJsonLd) }}
      />
      <PostDetail post={post} initialComments={comments} />
    </>
  );
}
