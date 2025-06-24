"use client";
import { List, Tag } from "antd";
import Image from "next/image";
import Link from "next/link";
import { EyeOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";
import type { LatestPostVO } from "../api/posts";
import SiteOwnerCard from "./SiteOwnerCard";

interface ArticleListLayoutProps {
  list: LatestPostVO[];
  total: number;
  siteOwner?: {
    name: string;
    avatar: string;
    bio: string;
    stats?: any;
  };
}

export default function ArticleList({ list, total, siteOwner }: ArticleListLayoutProps) {
  return (
    <div className="w-full max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-12 gap-8">
      {/* 左侧主内容区 8/12 */}
      <div className="md:col-span-8 flex flex-col gap-8 min-w-0">
        <section>
          <List
            itemLayout="horizontal"
            dataSource={list}
            locale={{ emptyText: <div className="py-8 text-center text-gray-400">暂无数据</div> }}
            renderItem={item => (
              <List.Item className="!p-4 !bg-white !rounded !shadow !overflow-hidden">
                <Link href={`/articles/${item.sug}`} className="w-full group grid grid-cols-6">
                  {/* 图片区域 grid 2/6，右侧加 pr-4 形成空隙 */}
                  <div className="col-span-2 h-32 relative flex items-center justify-center bg-gray-50 overflow-hidden pr-4">
                    {/* 标签区悬浮在图片左上角 */}
                    <div className="absolute top-2 left-2 flex flex-wrap gap-2 z-10">
                      {item.categories?.map((cat) => (
                        <Tag key={cat} color="blue">{cat}</Tag>
                      ))}
                      {item.tags?.map((tag) => (
                        <Tag key={tag} color="orange">#{tag}</Tag>
                      ))}
                    </div>
                    <Image src={item.cover_img} alt={item.title} fill className="object-cover" />
                  </div>
                  {/* 内容区域 grid 4/6 */}
                  <div className="col-span-4 flex flex-col justify-between p-4">
                    <div className="text-lg font-bold mb-1 group-hover:text-blue-600 transition-colors">{item.title}</div>
                    <div className="text-gray-700 mb-2 line-clamp-2 text-sm">{item.summary}</div>
                    <div className="flex items-center justify-between text-xs text-gray-400 mt-auto">
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
          <div className="text-right text-gray-400 mt-4">共 {total} 篇文章</div>
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
