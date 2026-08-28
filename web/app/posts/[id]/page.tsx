import PostDetail from "@/src/components/PostDetail";
import { getPostDetailOrNull } from "@/src/api/posts";
import { notFound } from "next/navigation";
import { getCommonConfig } from "@/src/api/config";
import type { Metadata } from "next";
import { resolvePublicUrl } from "@/src/utils/publicUrl";
import { getCommentsByPostId } from "@/src/api/comments";
import { buildBlogPostingJsonLd, serializeJsonLd } from "@/src/utils/seo";

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }): Promise<Metadata> {
  const { id } = await params;
  const post = await getPostDetailOrNull(id);
  if (!post) return {};
  const config = await getCommonConfig();
  return {
    title: `${post.title} - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: post.meta_description || post.summary,
    keywords: post.meta_keywords,
    authors: [{ name: post.author }],
    alternates: {
      canonical: `/posts/${id}`,
    },
    openGraph: {
      title: `${post.title} - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: post.meta_description || post.summary,
      url: process.env.BASE_HOST + `/posts/${id}`,
      images: post.cover_img ? [{ url: resolvePublicUrl(post.cover_img) }] : (config.seo_meta.og_image ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }] : undefined),
      siteName: config.website_meta.website_name,
      type: "article",
      publishedTime: new Date(post.created_at * 1000).toISOString(),
      modifiedTime: new Date((post.updated_at || post.created_at) * 1000).toISOString(),
      authors: [post.author],
      tags: post.tags.map((tag) => tag.name),
    },
    twitter: {
      card: "summary_large_image",
      title: `${post.title} - ${config.seo_meta.title || config.website_meta.website_name}`,
      description: post.meta_description || post.summary,
      images: post.cover_img
        ? [resolvePublicUrl(post.cover_img)!]
        : config.seo_meta.og_image
          ? [resolvePublicUrl(config.seo_meta.og_image)!]
          : undefined,
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
