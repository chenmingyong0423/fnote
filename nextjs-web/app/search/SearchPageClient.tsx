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
    <div className="w-full max-w-7xl mx-auto">
      <div className="mb-6 flex md:justify-start justify-center">
        <div className="w-full md:col-span-8" style={{ maxWidth: '66.6667%' }}>
          <Input.Search
            placeholder="请输入关键词..."
            allowClear
            enterButton="搜索"
            size="large"
            defaultValue={keyword}
            onSearch={handleSearch}
          />
        </div>
      </div>
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