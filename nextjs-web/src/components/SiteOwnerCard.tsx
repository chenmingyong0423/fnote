"use client";
import { Card, Avatar } from "antd";
import React from "react";
import SiteStats from "./SiteStats";

export interface SiteOwnerCardProps {
  name: string;
  avatar?: string;
  bio?: string;
  stats?: any;
}

export default function SiteOwnerCard({ name, avatar, bio, stats }: SiteOwnerCardProps) {
  return (
    <Card className="text-center p-6">
      <div className="flex flex-col items-center gap-3">
        {/* 头像 */}
        {avatar && avatar !== "" ? (
          <Avatar src={avatar} size={64} />
        ) : null}
        {/* 名字 */}
        <div className="font-bold text-lg dark:text-gray-200">{name}</div>
        {/* 简介 */}
        <div className="text-gray-500 text-sm dark:text-gray-400">{bio}</div>
        {/* 指标区 */}
        <SiteStats stats={stats} />
      </div>
    </Card>
  );
}