"use client";

import React, { useState } from "react";
import dynamic from "next/dynamic";
import {getPostList, PostListParams} from "../api/posts";
import { getCategoryNameStringByRoute } from "../api/category";
import { getTagNameByRoute } from "../api/tags";
import { getIndexConfig } from "../api/config";
import { getWebsiteStats } from "../api/stats";
import type { PostListResponse } from "../api/posts";
import {SiteOwnerCardProps} from "@/src/components/SiteOwnerCard";

const ArticleListNoSSR = dynamic(() => import("./ArticleList"), { ssr: false });

interface ArticleListContainerProps {
  category?: string;
  tag?: string;
  keyword?: string;
  field?: "latest" | "oldest" | "likes";
  page?: number;
  pageSize?: number;
}

export default function ArticleListContainer({
  category,
  tag,
  keyword,
  field = "latest",
  page = 1,
  pageSize = 10,
}: ArticleListContainerProps) {
  const [currentField, setCurrentField] = useState(field);
  const [currentPage, setCurrentPage] = useState(page);
  const [currentPageSize, setCurrentPageSize] = useState(pageSize);
  const [data, setData] = useState<PostListResponse>({
    list: [],
    PageNo: 1,
    PageSize: 10,
    totalPages: 0,
    totalCount: 0,
  });
  const [siteOwner, setSiteOwner] = useState<SiteOwnerCardProps>();

  React.useEffect(() => {
    async function fetchData() {
      let categoryName: string | undefined = undefined;
      let tagName: string | undefined = undefined;
      if (category) {
        categoryName = await getCategoryNameStringByRoute(category);
      }
      if (tag) {
        tagName = await getTagNameByRoute(tag);
      }
      
      // 监听 field/page/pageSize 变化都要刷新数据
      let sortField: string | undefined = undefined;
      let sortOrder: string | undefined = undefined;
      if (field === "latest") {
        sortField = "created_at";
        sortOrder = "DESC";
      } else if (field === "oldest") {
        sortField = "created_at";
        sortOrder = "ASC";
      } else if (field === "likes") {
        sortField = "like_count";
        sortOrder = "DESC";
      }
      const params: PostListParams = { pageNo: currentPage, pageSize: currentPageSize };
      if (sortField) params.sortField = sortField;
      if (sortOrder) params.sortOrder = sortOrder;
      if (categoryName) params.categories = [categoryName];
      if (tagName) params.tags = [tagName];
      if (keyword) params.keyword = keyword;
      const posts = await getPostList(params);
      setData(posts);
      const config = await getIndexConfig();
      const stats = await getWebsiteStats();
      setSiteOwner({
        name: config.website_config.website_owner,
        avatar: config.website_config.website_owner_avatar,
        bio: config.website_config.website_owner_profile,
        stats,
      });
    }
    fetchData().catch();
  }, [category, tag, keyword, field, currentField, currentPage, currentPageSize]);

  React.useEffect(() => {
    setCurrentField(field);
  }, [field]);
  React.useEffect(() => {
    setCurrentPage(page);
  }, [page]);
  React.useEffect(() => {
    setCurrentPageSize(pageSize);
  }, [pageSize]);

  return (
    <ArticleListNoSSR
      list={data.list}
      total={data.totalCount}
      siteOwner={siteOwner}
      field={currentField}
      currentPage={currentPage}
      pageSize={currentPageSize}
    />
  );
}
