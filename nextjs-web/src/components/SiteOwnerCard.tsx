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
        name: "",
        avatar: "",
        bio: ""
      };
  return (
    <Card className="text-center p-6">
      <div className="flex flex-col items-center gap-3">
        {/* 头像 */}
        <Avatar src={siteOwner.avatar} size={64} />
        {/* 名字 */}
        <div className="font-bold text-lg">{siteOwner.name}</div>
        {/* 简介 */}
        <div className="text-gray-500 text-sm">{siteOwner.bio}</div>
        {/* 指标区 */}
        <SiteStats />
      </div>
    </Card>
  );
}