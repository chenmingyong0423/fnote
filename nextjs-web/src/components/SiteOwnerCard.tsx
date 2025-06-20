"use client";
import { Card, Avatar } from "antd";
import React, { useEffect } from "react";
import { getWebsiteStats } from "../api/stats";
import { useConfigStore } from "../store/config";
import SiteStats from "./SiteStats";

export default function SiteOwnerCard() {
  const config = useConfigStore((s) => s.config);
  const setConfig = useConfigStore((s) => s.setConfig);

  useEffect(() => {
    if (config && !("stats" in config)) {
      getWebsiteStats().then((stats) => {
        setConfig({ ...config, stats });
      });
    }
  }, [config, setConfig]);

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
        {/* 指标区 */}
        <SiteStats />
      </div>
    </Card>
  );
}
