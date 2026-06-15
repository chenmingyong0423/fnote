"use client";

import { Avatar } from "antd";
import React from "react";
import SiteStats, { SiteStatsProps } from "./SiteStats";

export interface SiteOwnerCardProps {
  name: string;
  avatar?: string;
  bio?: string;
  stats?: SiteStatsProps;
  hasError?: boolean;
}

export default function SiteOwnerCard({
  name,
  avatar,
  bio,
  stats,
  hasError = false,
}: SiteOwnerCardProps) {
  return (
    <section className="overflow-hidden rounded-lg border border-gray-200/80 bg-white p-5 text-center shadow-sm dark:border-[#303030] dark:bg-[#141414] md:p-6">
      <div className="mx-auto mb-4 h-1 w-12 rounded-full bg-blue-500/80" />
      <div className="flex flex-col items-center gap-3">
        {avatar && avatar !== "" ? (
          <Avatar
            src={avatar}
            size={64}
            className="ring-4 ring-gray-100 dark:ring-[#232426]"
          />
        ) : null}
        <div className="break-words text-base font-bold text-gray-950 dark:text-gray-100 md:text-lg">
          {name}
        </div>
        <div className="max-w-[18rem] break-words text-sm leading-6 text-gray-500 dark:text-gray-400">
          {hasError ? "网站数据暂时异常" : bio}
        </div>
        <SiteStats stats={stats} />
      </div>
    </section>
  );
}
