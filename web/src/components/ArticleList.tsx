"use client";
import React, { useEffect, useState } from "react";
import { Pagination, Tag, Tabs, type PaginationProps } from "antd";
import Image from "next/image";
import Link from "next/link";
import { EyeOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";
import SiteOwnerCard, { SiteOwnerCardProps } from "./SiteOwnerCard";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import type { LatestPostVO } from "../api/posts";
import { formatDate } from "../utils/date";

interface ArticleListProps {
  list: LatestPostVO[];
  total: number;
  pageHeading?: string;
  siteOwner?: SiteOwnerCardProps;
  field?: "latest" | "oldest" | "likes";
  currentPage?: number;
  pageSize?: number;
  hasError?: boolean;
  hideSiteOwnerOnMobile?: boolean;
}

export default function ArticleList({
  list,
  total,
  pageHeading,
  siteOwner,
  field = "latest",
  currentPage = 1,
  pageSize = 10,
  hasError = false,
  hideSiteOwnerOnMobile = false,
}: ArticleListProps) {
  const [localField, setLocalField] = useState(field);
  const [localPage, setLocalPage] = useState(currentPage);
  const [localPageSize, setLocalPageSize] = useState(pageSize);
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const basePath = pathname.replace(/\/page\/[0-9]+$/, "");

  // 与点击跳转共用地址，确保服务端 HTML 中就有可抓取的分页链接。
  const getPageHref = (page: number, size = localPageSize, sort = localField) => {
    const params = new URLSearchParams(searchParams?.toString() ?? "");
    params.delete("page");
    if (size === 10) params.delete("pageSize");
    else params.set("pageSize", String(size));
    if (sort === "latest") params.delete("filter");
    else params.set("filter", sort);
    const query = params.toString();
    const path = page === 1 ? basePath : `${basePath}/page/${page}`;
    return query ? `${path}?${query}` : path;
  };

  useEffect(() => {
    setLocalField(field);
    setLocalPage(currentPage);
    setLocalPageSize(pageSize);
  }, [currentPage, field, pageSize]);

  // 排序切换
  const handleFilterChange = (value: string) => {
    setLocalField(value as "latest" | "oldest" | "likes");
    setLocalPage(1);
    router.push(getPageHref(1, localPageSize, value as typeof localField), { scroll: true });
  };

  // 分页切换
  const handlePageChange = (page: number, size: number) => {
    const targetPage = size === localPageSize ? page : 1;
    setLocalPage(targetPage);
    setLocalPageSize(size);
    router.push(getPageHref(targetPage, size), { scroll: true });
  };

  const renderPaginationItem: NonNullable<PaginationProps["itemRender"]> = (page, type, element) => {
    const totalPages = Math.ceil(total / localPageSize);
    if (page < 1 || page > totalPages ||
      (type === "prev" && localPage <= 1) ||
      (type === "next" && localPage >= totalPages)) return element;

    const original = React.isValidElement<{ className?: string; children?: React.ReactNode }>(element)
      ? element.props : undefined;
    const label = type === "prev" ? "上一页" : type === "next" ? "下一页" : `第 ${page} 页`;
    return (
      <a
        href={getPageHref(page)}
        className={original?.className}
        aria-label={label}
        aria-current={type === "page" && page === localPage ? "page" : undefined}
        onKeyDown={(event) => event.stopPropagation()}
        onClick={(event) => {
          // 避免父级 Pagination 重复导航，保留 Ctrl/Cmd 点击等浏览器行为。
          event.stopPropagation();
          if (event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
          event.preventDefault();
          handlePageChange(page, localPageSize);
        }}
      >
        {original?.children ?? label}
      </a>
    );
  };

  return (
    <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
      {pageHeading && (
        <h1 className="mb-5 text-xl font-bold text-gray-900 dark:text-gray-100 md:mb-6 md:text-2xl">
          {pageHeading}
        </h1>
      )}
      {/* 移动端：纵向布局，桌面端：网格布局 */}
      <div className="flex flex-col md:grid md:grid-cols-12 gap-6 md:gap-8 dark:text-gray-200">
        {/* 左侧主内容区 - 移动端全宽，桌面端 8/12 */}
        <div className={`w-full ${siteOwner ? "md:col-span-8" : "md:col-span-12"} flex flex-col gap-5 md:gap-8 min-w-0`}>
          <section>
            {/* 排序过滤选项 */}
            <div className="mb-4 overflow-x-auto overflow-y-hidden [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
              <Tabs
                activeKey={localField}
                onChange={handleFilterChange}
                size="small"
                className="min-w-max"
                items={[
                  { key: "latest", label: "最新发布" },
                  { key: "oldest", label: "最早发布" },
                  { key: "likes", label: "点赞最多" },
                ]}
              />
            </div>
          {list.length > 0 ? (
            <div className="flex flex-col gap-4">
              {list.map((item) => (
                <div
                  key={item.sug}
                  className="glass-card w-full rounded-lg overflow-hidden p-3 transition-transform duration-200 group/article md:p-4 md:hover:-translate-y-2 relative"
                >
                  {/* 下划线动画，item hover 时从中间向两边展开 */}
                  <span className="pointer-events-none absolute left-1/2 bottom-0 w-0 h-0.5 bg-blue-500 rounded-full transition-all duration-300 group-hover/article:w-full group-hover/article:left-0"></span>
                  <Link
                    href={`/posts/${item.sug}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="w-full group block"
                  >
                    {/* 移动端：纵向布局，桌面端：横向布局 */}
                    <div className="flex flex-col gap-4 md:grid md:grid-cols-6 md:gap-5">
                      {/* 图片区域 - 移动端全宽，桌面端 2/6 */}
                      <div className="w-full md:col-span-2 h-40 sm:h-48 md:h-32 relative flex items-center justify-center overflow-hidden rounded-md border border-white/60 bg-white/35 shadow-sm dark:border-slate-500/20 dark:bg-slate-950/30">
                        {/* 标签区悬浮在图片左上角，初始隐藏，hover 时滑入 */}
                        <div className="absolute top-2 left-2 right-2 flex flex-wrap gap-1.5 md:gap-2 z-10 transition-all duration-300 md:-translate-x-6 md:opacity-0 md:group-hover/article:translate-x-0 md:group-hover/article:opacity-100">
                          {item.categories?.map((cat, index) => (
                            <Tag key={`category-${cat}-${index}`} color="#2DB7F5" style={{ color: '#fff', border: 'none', fontSize: '12px' }}>{cat}</Tag>
                          ))}
                          {item.tags?.map((tag, index) => (
                            <Tag key={`tag-${tag}-${index}`} color="#FB923C" style={{ color: '#fff', border: 'none', fontSize: '12px' }}>#{tag}</Tag>
                          ))}
                        </div>
                        <Image src={item.cover_img} alt={item.title} fill className="object-contain" />
                      </div>
                      {/* 内容区域 - 移动端全宽，桌面端 4/6 */}
                      <div className="w-full md:col-span-4 flex flex-col justify-between md:py-1 relative">
                        <h2 className="text-base md:text-lg font-bold mb-2 group-hover:text-blue-600 transition-colors dark:text-gray-200 dark:group-hover:text-gray-100 line-clamp-2">{item.title}</h2>
                        <div className="text-gray-700 mb-3 line-clamp-3 text-sm dark:text-gray-400 flex-1">{item.summary}</div>
                        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between text-xs text-gray-400 mt-auto dark:text-gray-500">
                          <div className="flex flex-wrap items-center gap-3 md:gap-4">
                            <span className="flex items-center gap-1"><EyeOutlined /> {item.visit_count}</span>
                            <span className="flex items-center gap-1"><LikeOutlined /> {item.like_count}</span>
                            <span className="flex items-center gap-1"><MessageOutlined /> {item.comment_count}</span>
                          </div>
                          <div className="sm:text-right whitespace-nowrap text-xs">{formatDate(item.created_at)}</div>
                        </div>
                      </div>
                    </div>
                  </Link>
                </div>
              ))}
            </div>
          ) : (
            <div className="py-8 text-center text-gray-400">
              {hasError ? "网站数据暂时异常" : "暂无数据"}
            </div>
          )}
          <div className="article-pagination glass-surface mt-4 overflow-x-auto rounded-lg px-3 py-2">
            <Pagination
              current={localPage}
              pageSize={localPageSize}
              total={total}
              onChange={handlePageChange}
              itemRender={renderPaginationItem}
              showSizeChanger={{
                getPopupContainer: () => document.body,
                classNames: {
                  popup: {
                    root: "article-pagination-size-dropdown",
                  },
                },
              }}
              pageSizeOptions={["5", "10", "20", "50"]}
              showQuickJumper={true}
              showTotal={total => `共 ${total} 篇文章`}
              responsive={true}
              simple={false}
            />
          </div>
        </section>
      </div>
      {/* 右侧信息区 - 移动端全宽，桌面端 4/12 */}
      {siteOwner && (
        <div className={`w-full md:col-span-4 ${hideSiteOwnerOnMobile ? "hidden md:flex" : "flex"} flex-col gap-6 md:gap-8 min-w-0`}>
          <SiteOwnerCard
            name={siteOwner.name}
            avatar={siteOwner.avatar}
            bio={siteOwner.bio}
            stats={siteOwner.stats}
            hasError={siteOwner.hasError}
          />
        </div>
      )}
    </div>
  </div>
  );
}
