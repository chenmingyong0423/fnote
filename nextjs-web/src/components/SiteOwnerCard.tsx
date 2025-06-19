"use client";
import { Card, Avatar } from "antd";
import React from "react";
import { useConfigStore } from "../store/config";

export default function SiteOwnerCard() {
  const config = useConfigStore((s) => s.config);
  const siteOwner = config
    ? {
        name: config.website_config.website_owner,
        avatar: config.website_config.website_owner_avatar,
        bio: config.website_config.website_owner_profile,
        links: [],
      }
    : {
        name: "站长小明",
        avatar: "/logo.png",
        bio: "全栈开发者，热爱开源与分享。专注于前端、Node.js、云原生。",
        links: [
          { label: "GitHub", url: "https://github.com/owner" },
          { label: "博客", url: "/about" },
        ],
      };
  return (
    <Card className="text-center p-6">
      <div className="flex flex-col items-center gap-6">
        {/* 头像 */}
        <Avatar src={siteOwner.avatar} size={64} className="mb-3" />
        {/* 名字 */}
        <div className="font-bold text-lg mb-3">{siteOwner.name}</div>
        {/* 简介 */}
        <div className="text-gray-500 mb-5 text-sm">{siteOwner.bio}</div>
        {/* 指标区，两行三列布局 */}
        <div className="w-full border-t border-b border-gray-200 py-3 grid grid-cols-3 gap-y-2 text-xs text-gray-600 divide-x divide-gray-200">
          <div className="col-span-1 flex flex-col items-center">
            <span>文章</span>
            <span className="font-bold">66</span>
          </div>
          <div className="col-span-1 flex flex-col items-center">
            <span>分类</span>
            <span className="font-bold">8</span>
          </div>
          <div className="col-span-1 flex flex-col items-center">
            <span>标签</span>
            <span className="font-bold">20</span>
          </div>
          <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
            <span>评论</span>
            <span className="font-bold">123</span>
          </div>
          <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
            <span>点赞</span>
            <span className="font-bold">888</span>
          </div>
          <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
            <span>浏览</span>
            <span className="font-bold">9999</span>
          </div>
        </div>
      </div>
    </Card>
  );
}
