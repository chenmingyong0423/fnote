import React from "react";
import { getPostList } from "../api/posts";
import { getCategoryNameStringByRoute } from "../api/category";
import { getTagNameByRoute } from "../api/tags";
import { getIndexConfig } from "../api/config";
import { getWebsiteStats } from "../api/stats";
import ArticleList from "./ArticleList";

interface ArticleListContainerProps {
  category?: string;
  tag?: string;
  keyword?: string;
  field?: "latest" | "oldest" | "likes";
}

export default async function ArticleListContainer({
  category,
  tag,
  keyword,
  field = "latest",
}: ArticleListContainerProps) {
  let categoryName: string | undefined = undefined;
  let tagName: string | undefined = undefined;
  if (category) {
    categoryName = await getCategoryNameStringByRoute(category);
  }
  if (tag) {
    tagName = await getTagNameByRoute(tag);
  }
  const params: any = { pageNo: 1, pageSize: 10, field };
  if (categoryName) params.categories = [categoryName];
  if (tagName) params.tags = [tagName];
  if (keyword) params.keyword = keyword;

  const data = await getPostList(params);
  const config = await getIndexConfig();
  const stats = await getWebsiteStats();

  return <ArticleList list={data.list} total={data.totalCount} siteOwner={{
    name: config.website_config.website_owner,
    avatar: config.website_config.website_owner_avatar,
    bio: config.website_config.website_owner_profile,
    stats,
  }} />;
}
