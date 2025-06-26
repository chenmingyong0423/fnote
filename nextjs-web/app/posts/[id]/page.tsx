import React from "react";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import rehypeHighlight from "rehype-highlight";
import "highlight.js/styles/github-dark.css";
import { extractToc, Toc } from "@/src/components/Toc";
import { PostActions } from "@/src/components/PostActions";

function genHeadingId(text: string) {
  return text.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, "-").replace(/^-+|-+$/g, "");
}

export default async function PostDetailPage({ params }: { params: { id: string } }) {
  let post;
  try {
    const res = await getPostDetail(params.id);
    if (res.code !== 0 || !res.data) return notFound();
    post = res.data;
  } catch {
    return notFound();
  }

  const toc = extractToc(post.content);

  return (
    <div className="w-full max-w-7xl mx-auto flex flex-col lg:flex-row gap-8 bg-white dark:bg-[#141414] rounded-xl shadow-sm p-6 mt-8 mb-12">
      <div className="flex-1 min-w-0">
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
            components={{
              h1: ({node, children, ...props}) => <h1 id={genHeadingId(String(children))} {...props}>{children}</h1>,
              h2: ({node, children, ...props}) => <h2 id={genHeadingId(String(children))} {...props}>{children}</h2>,
              h3: ({node, children, ...props}) => <h3 id={genHeadingId(String(children))} {...props}>{children}</h3>,
              h4: ({node, children, ...props}) => <h4 id={genHeadingId(String(children))} {...props}>{children}</h4>,
              h5: ({node, children, ...props}) => <h5 id={genHeadingId(String(children))} {...props}>{children}</h5>,
              h6: ({node, children, ...props}) => <h6 id={genHeadingId(String(children))} {...props}>{children}</h6>,
            }}
          >
            {post.content}
          </ReactMarkdown>
        </article>
      </div>
      {/* 右侧悬浮区：操作区在上，目录在下，整体 sticky，避免重叠 */}
      <div className="w-full lg:w-64 flex-shrink-0">
        <div className="sticky top-24 flex flex-col gap-4">
          <PostActions />
          <Toc toc={toc} />
        </div>
      </div>
    </div>
  );
}
