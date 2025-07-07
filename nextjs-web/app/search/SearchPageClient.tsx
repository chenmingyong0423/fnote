"use client";
import React, {Suspense} from "react";
import { Input } from "antd";
import { useRouter, useSearchParams } from "next/navigation";
import ArticleList from "@/src/components/ArticleList";

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
  list: any[];
  total: number;
  siteOwner: any;
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
    router.replace("?" + newParams.toString());
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
  list: any[];
  total: number;
  siteOwner: any;
}) {
  return (
      <Suspense>
        <SearchPage keyword={keyword} field={field} page={page} pageSize={pageSize} list={list} total={total} siteOwner={siteOwner} />
      </Suspense>
  )
}