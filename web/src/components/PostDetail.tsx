import React from "react";
import { MarkdownPreview } from "@/src/components/MarkdownPreview";
import { extractToc, Toc } from "@/src/components/Toc";
import { PostActions } from "@/src/components/PostActions";
import { Comments } from "@/src/components/Comments";
import PostSeoClient from "@/src/components/PostSeoClient";
import { formatDate } from "@/src/utils/date";
import type { CommentItem } from "@/src/api/comments";

interface PostDetailProps {
  post: import("@/src/api/posts").PostDetail;
  initialComments?: CommentItem[];
}

const stripMarkdown = (content: string) => {
  return content
    .replace(/```[\s\S]*?```/g, " ")
    .replace(/`([^`]+)`/g, "$1")
    .replace(/!\[[^\]]*]\([^)]*\)/g, " ")
    .replace(/\[([^\]]+)]\([^)]*\)/g, "$1")
    .replace(/<[^>]+>/g, " ")
    .replace(/^#{1,6}\s+/gm, "")
    .replace(/^>\s?/gm, "")
    .replace(/^[\s>*+-]*\d+\.\s+/gm, "")
    .replace(/^[\s>*+-]+/gm, "")
    .replace(/[*_~>#|[\]()]/g, " ")
    .replace(/\s+/g, " ")
    .trim();
};

const countReadableWords = (content: string) => {
  const chineseCount = content.match(/[\u4e00-\u9fa5]/g)?.length || 0;
  const latinContent = content.replace(/[\u4e00-\u9fa5]/g, " ");
  const latinWordCount =
    latinContent.match(/[a-zA-Z0-9]+(?:[-'][a-zA-Z0-9]+)*/g)?.length || 0;
  return chineseCount + latinWordCount;
};

const getDisplayWordCount = (post: import("@/src/api/posts").PostDetail) => {
  if (post.word_count > 0) {
    return post.word_count;
  }
  return countReadableWords(stripMarkdown(post.content || ""));
};

const getReadingMinutes = (wordCount: number) => {
  if (wordCount === 0) {
    return 0;
  }
  return Math.max(1, Math.ceil(wordCount / 400));
};

const PostDetail: React.FC<PostDetailProps> = ({ post, initialComments }) => {
  const toc = extractToc(post.content);
  const wordCount = getDisplayWordCount(post);
  const readingMinutes = getReadingMinutes(wordCount);

  return (
    <>
      <PostSeoClient
        title={post.title}
        description={post.meta_description || post.summary || ""}
        keywords={post.meta_keywords || post.tags?.map((t) => t.name).join(",")}
        coverImg={post.cover_img}
      />
      <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
        <div className="glass-surface flex flex-col lg:flex-row gap-5 lg:gap-8 rounded-xl p-4 md:p-6 mb-8 md:mb-12">
          <div className="flex-1 min-w-0">
            <h1 className="text-xl leading-snug md:text-3xl font-bold mb-4 dark:text-gray-100">
              {post.title}
            </h1>
            <div className="flex flex-wrap gap-x-3 gap-y-2 md:gap-4 text-xs md:text-sm text-gray-500 dark:text-gray-400 mb-5 md:mb-6">
              <span>作者：{post.author}</span>
              <span>分类：{post.category.map((c) => c.name).join(", ")}</span>
              <span className="hidden sm:inline">
                标签：{post.tags.map((t) => t.name).join(", ")}
              </span>
              <span>发布：{formatDate(post.created_at)}</span>
              {wordCount > 0 && (
                <>
                  <span>字数：{wordCount}</span>
                  <span>阅读：约 {readingMinutes} 分钟</span>
                </>
              )}
              <span className="flex flex-wrap items-center gap-2 md:gap-3">
                <span>浏览：{post.visit_count}</span>
              </span>
            </div>
            {/* 移动端显示标签 */}
            <div className="sm:hidden mb-4 text-xs text-gray-500 dark:text-gray-400">
              标签：{post.tags.map((t) => t.name).join(", ")}
            </div>
            <article className="glass-card overflow-hidden rounded-lg px-4 py-5 md:px-7 md:py-6 prose prose-sm md:prose-lg max-w-none dark:prose-invert">
              <MarkdownPreview
                content={post.content}
                theme="blog"
                signatureText={post.author}
              />
            </article>
            {/* 版权信息区 */}
            <div className="mt-8 md:mt-10 p-3 md:p-4 rounded-lg border border-white/70 bg-white/45 text-xs md:text-sm text-gray-600 backdrop-blur dark:border-white/10 dark:bg-slate-900/35 dark:text-gray-400">
              <div className="break-all">
                本文链接：
                <span>{`${process.env.BASE_HOST || ""}/posts/${post._id}`}</span>
              </div>
              <div className="mt-2">
                版权声明：本文由{" "}
                <span className="font-semibold">{post.author}</span>{" "}
                原创发布，如需转载请遵循
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
          {/* 右侧悬浮区：操作区在上，目录在下，整体 sticky，避免重叠。移动端隐藏 */}
          <div className="w-80 hidden lg:flex flex-col gap-4 sticky top-24 h-fit">
            <PostActions
              postId={post._id}
              isLiked={post.is_liked}
              likeCount={post.like_count}
              commentCount={post.comment_count}
            />
            <Toc toc={toc} />
          </div>
        </div>
      </div>
      {/* 移动端操作按钮，固定在底部 */}
      <div className="lg:hidden fixed bottom-3 right-3 z-50 max-w-[calc(100vw-1.5rem)]">
        <PostActions
          postId={post._id}
          isLiked={post.is_liked}
          likeCount={post.like_count}
          commentCount={post.comment_count}
        />
      </div>
      <div className="w-full max-w-7xl mx-auto px-4 md:px-0 mt-0 mb-10 md:mb-12">
        <Comments postId={post._id} initialComments={initialComments} />
      </div>
    </>
  );
};

export default PostDetail;
