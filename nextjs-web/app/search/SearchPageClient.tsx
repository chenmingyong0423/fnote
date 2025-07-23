"use client";
import React, {Suspense} from "react";
import { Input } from "antd";
import { useRouter, useSearchParams } from "next/navigation";
import ArticleList from "@/src/components/ArticleList";
import {LatestPostVO} from "@/src/api/posts";
import {SiteOwnerCardProps} from "@/src/components/SiteOwnerCard";

function SearchPage({
  keyword,
  field,
  page,
  pageSize,
  list,
  total,
  siteOwner,
}: {
  keyword: string;
  field: "latest" | "oldest" | "likes";
  page: number;
  pageSize: number;
  list: LatestPostVO[];
  total: number;
  siteOwner: SiteOwnerCardProps;
}) {
  const router = useRouter();
  const params = useSearchParams();
  
  const handleSearch = (value: string) => {
    const newParams = new URLSearchParams(params?.toString());
    if (value) {
      newParams.set("keyword", value);
    } else {
      newParams.delete("keyword");
    }
    newParams.delete("page"); // 搜索时重置分页
    router.replace("/search?" + newParams.toString());
  };

  // 搜索页面专用的分页处理函数
  const handlePageChange = (targetPage: number, size: number, currentField: string) => {
    const newParams = new URLSearchParams(params?.toString());
    newParams.set("pageSize", String(size));
    newParams.set("filter", currentField);
    
    // 搜索页面的分页跳转逻辑
    if (targetPage === 1) {
      router.push(`/search?${newParams.toString()}`);
    } else {
      router.push(`/search/page/${targetPage}?${newParams.toString()}`);
    }
  };

  return (
    <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
      {/* 搜索框区域 */}
      <div className="mb-6 bg-white dark:bg-[#141414] rounded-lg shadow-sm p-4 md:p-6">
        <h1 className="text-xl md:text-2xl font-bold mb-4 dark:text-gray-100">搜索文章</h1>
        <div className="w-full">
          <Input.Search
            placeholder="请输入关键词搜索文章..."
            allowClear
            enterButton="搜索"
            size="large"
            defaultValue={keyword}
            onSearch={handleSearch}
            className="w-full"
          />
        </div>
        {keyword && (
          <div className="mt-3 text-sm text-gray-500 dark:text-gray-400">
            搜索关键词：<span className="font-medium text-blue-600 dark:text-blue-400">"{keyword}"</span>
            {total > 0 && <span className="ml-2">找到 {total} 篇相关文章</span>}
          </div>
        )}
      </div>
      
      {/* 文章列表 */}
      <ArticleList
        list={list}
        total={total}
        siteOwner={siteOwner}
        field={field}
        currentPage={page}
        pageSize={pageSize}
        onPageChange={handlePageChange}
      />
    </div>
  );
}

export default function SearchPageClient({
  keyword,
  field,
  page,
  pageSize,
  list,
  total,
  siteOwner,
}: {
  keyword: string;
  field: "latest" | "oldest" | "likes";
  page: number;
  pageSize: number;
  list: LatestPostVO[];
  total: number;
  siteOwner: SiteOwnerCardProps;
}) {
  return (
      <Suspense>
        <SearchPage keyword={keyword} field={field} page={page} pageSize={pageSize} list={list} total={total} siteOwner={siteOwner} />
      </Suspense>
  )
}