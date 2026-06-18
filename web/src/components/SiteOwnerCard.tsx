"use client";

import { Avatar } from "antd";
import React from "react";
import SiteStats, { SiteStatsProps } from "./SiteStats";
import SocialLinks from "./SocialLinks";
import type { SocialInfoVO } from "../api/config";

export interface SiteOwnerCardProps {
  name: string;
  avatar?: string;
  bio?: string;
  socialInfo?: SocialInfoVO[];
  stats?: SiteStatsProps;
  hasError?: boolean;
}

export default function SiteOwnerCard({
  name,
  avatar,
  bio,
  socialInfo,
  stats,
  hasError = false,
}: SiteOwnerCardProps) {
  return (
    <section className="glass-surface overflow-hidden rounded-lg p-5 text-center md:p-6">
      <div className="mx-auto mb-4 h-1 w-12 rounded-full bg-blue-500/80" />
      <div className="flex flex-col items-center gap-3">
        {avatar && avatar !== "" ? (
          <span className="inline-flex rounded-full transition-transform duration-700 hover:rotate-[360deg]">
            <Avatar
              src={avatar}
              size={64}
              className="ring-4 ring-gray-100 dark:ring-[#232426]"
            />
          </span>
        ) : null}
        <div className="break-words text-base font-bold text-gray-950 dark:text-gray-100 md:text-lg">
          {name}
        </div>
        <div className="max-w-[18rem] break-words text-sm leading-6 text-gray-500 dark:text-gray-400">
          {hasError ? "网站数据暂时异常" : bio}
        </div>
        {!hasError && <SocialLinks items={socialInfo} />}
        <SiteStats stats={stats} />
      </div>
    </section>
  );
}
