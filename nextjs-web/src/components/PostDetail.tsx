import React from "react";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";
import { MarkdownPreview, genHeadingId } from "@/src/components/MarkdownPreview";
import { extractToc, Toc } from "@/src/components/Toc";
import { PostActions } from "@/src/components/PostActions";
import { Comments } from "@/src/components/Comments";
import { log } from "console";
import PostSeoClient from "@/src/components/PostSeoClient";

interface PostDetailProps {
  id: string;
}

const PostDetail: React.FC<PostDetailProps> = async ({ id }) => {
  let post;
  try {
    const res = await getPostDetail(id);
    if (res.code !== 0 || !res.data) return notFound();
    post = res.data;
    // 推送百度索引（仅在服务端渲染时调用，防止多次推送）
    if (typeof window === "undefined" && post._id) {
        log(`Pushing post ${post._id} to Baidu index...`);
      // 动态 import，避免打包体积增加
      const { baiduPushIndex } = await import("@/src/api/baiduPush");
      const url = `${process.env.BASE_HOST || ''}/posts/${post._id}`;
      try {
        await baiduPushIndex({ urls: url });
      } catch {}
    }
  } catch {
    return notFound();
  }

  const toc = extractToc(post.content);

  return (
    <>
      <PostSeoClient
        title={post.title}
        description={post.meta_description || post.summary || ''}
        keywords={post.meta_keywords || post.tags?.map(t => t.name).join(",")}
        coverImg={post.cover_img}
      />
      <div className="w-full max-w-7xl mx-auto flex flex-col lg:flex-row gap-8 bg-white dark:bg-[#141414] rounded-xl shadow-sm p-6 mb-12">
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
          <article>
            <MarkdownPreview content={post.content} />
          </article>
          {/* 版权信息区 */}
          <div className="mt-10 p-4 rounded-lg bg-gray-50 dark:bg-[#232426] border border-gray-100 dark:border-gray-700 text-sm text-gray-600 dark:text-gray-400">
            <div>
              本文链接：<span className="break-all">{`${process.env.BASE_HOST || ''}/posts/${post._id}`}</span>
            </div>
            <div className="mt-2">
              版权声明：本文由 <span className="font-semibold">{post.author}</span> 原创发布，如需转载请遵循
              <a
                href="https://creativecommons.org/licenses/by-nc-sa/4.0/deed.en"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 dark:text-blue-400 underline mx-1"
              >
                署名-非商业性使用-相同方式共享 4.0 国际 (CC BY-NC-SA 4.0)
              </a>
              许可协议授权
            </div>
          </div>
        </div>
        {/* 右侧悬浮区：操作区在上，目录在下，整体 sticky，避免重叠 */}
        <div className="w-80 hidden lg:flex flex-col gap-4 sticky top-24 h-fit">
          <PostActions postId={post._id} />
          <Toc toc={toc} />
        </div>
      </div>
      <div className="w-full max-w-7xl mx-auto mt-0 mb-12">
        <Comments postId={post._id} />
      </div>
    </>
  );
};

export default PostDetail;
