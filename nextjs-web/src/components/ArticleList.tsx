"use client";
import React, { useState } from "react";
import { Pagination, List, Tag, Tabs } from "antd";
import Image from "next/image";
import Link from "next/link";
import { EyeOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";
import SiteOwnerCard, { SiteOwnerCardProps } from "./SiteOwnerCard";
import { useRouter, useSearchParams } from "next/navigation";
import type { LatestPostVO } from "../api/posts";

interface ArticleListProps {
  list: LatestPostVO[];
  total: number;
  siteOwner?: SiteOwnerCardProps;
  field?: "latest" | "oldest" | "likes";
  currentPage?: number;
  pageSize?: number;
  onPageChange?: (page: number, pageSize: number, field: string) => void;
}

export default function ArticleList({
  list,
  total,
  siteOwner,
  field = "latest",
  currentPage = 1,
  pageSize = 10,
  onPageChange,
}: ArticleListProps) {
  const [localField, setLocalField] = useState(field);
  const [localPage, setLocalPage] = useState(currentPage);
  const [localPageSize, setLocalPageSize] = useState(pageSize);
  const router = useRouter();
  const searchParams = useSearchParams();

  // 排序切换
  const handleFilterChange = (value: string) => {
    setLocalField(value as "latest" | "oldest" | "likes");
    setLocalPage(1);
    if (onPageChange) onPageChange(1, localPageSize, value);
    // 跳回第一页
    const params = new URLSearchParams(searchParams?.toString() ?? "");
    params.set("filter", value);
    params.delete("page");
    const base = window.location.pathname.replace(/\/page\/[0-9]+$/, "");
    router.push(`${base}?${params.toString()}`);
  };

  // 分页切换
  const handlePageChange = (page: number, size: number) => {
    setLocalPage(page);
    setLocalPageSize(size);
    if (onPageChange) onPageChange(page, size, localField);
    const params = new URLSearchParams(searchParams?.toString() ?? "");
    params.set("pageSize", String(size));
    const base = window.location.pathname.replace(/\/page\/[0-9]+$/, "");
    const targetPage = page === 1 ? "" : `/page/${page}`;
    router.push(`${base}${targetPage}?${params.toString()}`);
  };

  return (
    <div className="w-full max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-12 gap-8 dark:text-gray-200">
      {/* 左侧主内容区 8/12 */}
      <div className="md:col-span-8 flex flex-col gap-8 min-w-0">
        <section>
          {/* 排序过滤选项 */}
          <div className="flex items-center justify-between mb-4">
            <Tabs
              activeKey={localField}
              onChange={handleFilterChange}
              items={[
                { key: "latest", label: "最新发布" },
                { key: "oldest", label: "最早发布" },
                { key: "likes", label: "点赞最多" },
              ]}
            />
          </div>
          <List
            itemLayout="horizontal"
            dataSource={list}
            locale={{ emptyText: <div className="py-8 text-center text-gray-400">暂无数据</div> }}
            renderItem={item => (
              <List.Item className="!p-4 !bg-white dark:!bg-[#141414] dark:border dark:border-[#303030] !rounded !shadow !overflow-hidden my-4 transition-transform duration-200 group/article hover:-translate-y-2 relative">
                {/* 下划线动画，item hover 时从中间向两边展开 */}
                <span className="pointer-events-none absolute left-1/2 bottom-0 w-0 h-0.5 bg-blue-500 rounded-full transition-all duration-300 group-hover/article:w-full group-hover/article:left-0"></span>
                <Link href={`/posts/${item.sug}`} className="w-full group grid grid-cols-6">
                  {/* 图片区域 grid 2/6，右侧加 pr-4 形成空隙 */}
                  <div className="col-span-2 h-32 relative flex items-center justify-center bg-gray-50 overflow-hidden pr-4">
                    {/* 标签区悬浮在图片左上角，初始隐藏，hover 时滑入 */}
                    <div className="absolute top-2 left-2 flex flex-wrap gap-2 z-10 transition-all duration-300 -translate-x-6 opacity-0 group-hover/article:translate-x-0 group-hover/article:opacity-100">
                      {item.categories?.map((cat) => (
                        <Tag key={cat} color="#2DB7F5" style={{ color: '#fff', border: 'none' }}>{cat}</Tag>
                      ))}
                      {item.tags?.map((tag) => (
                        <Tag key={tag} color="#FB923C" style={{ color: '#fff', border: 'none' }}>#{tag}</Tag>
                      ))}
                    </div>
                    <Image src={item.cover_img} alt={item.title} fill className="object-cover" />
                  </div>
                  {/* 内容区域 grid 4/6 */}
                  <div className="col-span-4 flex flex-col justify-between p-4 relative">
                    <div className="text-lg font-bold mb-1 group-hover:text-blue-600 transition-colors dark:text-gray-200 dark:group-hover:text-gray-100">{item.title}</div>
                    <div className="text-gray-700 mb-2 line-clamp-2 text-sm dark:text-gray-400">{item.summary}</div>
                    <div className="flex items-center justify-between text-xs text-gray-400 mt-auto dark:text-gray-500">
                      <div className="flex items-center gap-4">
                        <span className="flex items-center gap-1"><EyeOutlined /> {item.visit_count}</span>
                        <span className="flex items-center gap-1"><LikeOutlined /> {item.like_count}</span>
                        <span className="flex items-center gap-1"><MessageOutlined /> {item.comment_count}</span>
                      </div>
                      <div className="text-right whitespace-nowrap">{new Date(item.created_at * 1000).toLocaleDateString()}</div>
                    </div>
                  </div>
                </Link>
              </List.Item>
            )}
          />
          <div className="flex justify-end mt-4">
            <Pagination
              current={localPage}
              pageSize={localPageSize}
              total={total}
              onChange={handlePageChange}
              showSizeChanger={true}
              pageSizeOptions={["5", "10", "20", "50"]}
              onShowSizeChange={(_, size) => handlePageChange(1, size)}
              showQuickJumper={true}
              showTotal={total => `共 ${total} 篇文章`}
            />
          </div>
        </section>
      </div>
      {/* 右侧信息区 4/12 */}
      <div className="md:col-span-4 flex flex-col gap-8 min-w-0">
        {siteOwner && (
          <SiteOwnerCard
            name={siteOwner.name}
            avatar={siteOwner.avatar}
            bio={siteOwner.bio}
            stats={siteOwner.stats}
          />
        )}
      </div>
    </div>
  );
}
