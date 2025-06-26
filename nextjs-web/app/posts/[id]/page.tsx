import React from "react";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import rehypeHighlight from "rehype-highlight";
import "highlight.js/styles/github-dark.css";

export default async function PostDetailPage({ params }: { params: { id: string } }) {
  let post;
  try {
    const res = await getPostDetail(params.id);
    if (res.code !== 0 || !res.data) return notFound();
    post = res.data;
  } catch {
    return notFound();
  }

  return (
    <div className="w-full max-w-7xl mx-auto bg-white dark:bg-[#141414] rounded-xl shadow-sm p-6 mt-8 mb-12">
      <h1 className="text-3xl font-bold mb-4 dark:text-gray-100">{post.title}</h1>
      <div className="flex flex-wrap gap-4 text-sm text-gray-500 dark:text-gray-400 mb-6">
        <span>作者：{post.author}</span>
        <span>分类：{post.category.map(c => c.name).join(', ')}</span>
        <span>标签：{post.tags.map(t => t.name).join(', ')}</span>
        <span>发布时间：{new Date(post.created_at * 1000).toLocaleString()}</span>
        <span>浏览：{post.visit_count}</span>
        <span>评论：{post.comment_count}</span>
        <span>点赞：{post.like_count}</span>
      </div>
      {post.cover_img && (
        <img src={post.cover_img} alt="cover" className="w-full max-h-96 object-cover rounded-lg mb-6" />
      )}
      <article className="prose prose-neutral dark:prose-invert max-w-none">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          rehypePlugins={[rehypeRaw, rehypeHighlight]}
        >
          {post.content}
        </ReactMarkdown>
      </article>
    </div>
  );
}
