"use client";

import { CloseOutlined } from "@ant-design/icons";
import { Avatar, Button } from "antd";
import React, { useEffect, useRef, useState } from "react";
import LatestComments, { type LatestComment } from "./LatestComments";
import SiteOwnerCard, { type SiteOwnerCardProps } from "./SiteOwnerCard";

interface HomeSidebarProps {
  siteOwner: SiteOwnerCardProps;
  comments: LatestComment[];
  hasCommentError?: boolean;
}

function useMediaQuery(query: string) {
  const [matches, setMatches] = useState<boolean | null>(null);

  useEffect(() => {
    const media = window.matchMedia(query);
    const update = () => setMatches(media.matches);

    update();
    media.addEventListener("change", update);
    return () => media.removeEventListener("change", update);
  }, [query]);

  return matches;
}

export default function HomeSidebar({
  siteOwner,
  comments,
  hasCommentError = false,
}: HomeSidebarProps) {
  const isDesktop = useMediaQuery("(min-width: 768px)");
  const [open, setOpen] = useState(false);
  const [drawerMounted, setDrawerMounted] = useState(false);
  const closeTimerRef = useRef<number | null>(null);

  useEffect(() => {
    if (!drawerMounted) return;

    const originalOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = originalOverflow;
    };
  }, [drawerMounted]);

  useEffect(() => {
    if (!isDesktop) return;

    setOpen(false);
    setDrawerMounted(false);
  }, [isDesktop]);

  useEffect(() => {
    return () => {
      if (closeTimerRef.current) {
        window.clearTimeout(closeTimerRef.current);
      }
    };
  }, []);

  const openOwnerPanel = () => {
    if (closeTimerRef.current) {
      window.clearTimeout(closeTimerRef.current);
      closeTimerRef.current = null;
    }

    setDrawerMounted(true);
    window.requestAnimationFrame(() => {
      setOpen(true);
    });
  };

  const closeOwnerPanel = () => {
    setOpen(false);
    closeTimerRef.current = window.setTimeout(() => {
      setDrawerMounted(false);
      closeTimerRef.current = null;
    }, 320);
  };

  if (isDesktop === null) {
    return <div className="hidden md:flex md:col-span-4" />;
  }

  if (isDesktop) {
    return (
      <div className="w-full md:col-span-4 flex flex-col gap-6 md:gap-8 min-w-0">
        <SiteOwnerCard {...siteOwner} />
        <LatestComments comments={comments} hasError={hasCommentError} />
      </div>
    );
  }

  const statsItems = siteOwner.stats
    ? [
        { label: "文章", value: siteOwner.stats.post_count },
        { label: "分类", value: siteOwner.stats.category_count },
        { label: "标签", value: siteOwner.stats.tag_count },
        { label: "评论", value: siteOwner.stats.comment_count },
        { label: "点赞", value: siteOwner.stats.like_count },
        { label: "浏览", value: siteOwner.stats.website_view_count },
      ]
    : [];

  return (
    <>
      <button
        type="button"
        aria-label="打开站长信息"
        onClick={openOwnerPanel}
        className="fixed bottom-5 right-4 z-40 flex h-14 w-14 items-center justify-center rounded-full border border-white/80 bg-white shadow-lg shadow-black/15 transition-transform active:scale-95 dark:border-gray-700 dark:bg-[#232426]"
      >
        {siteOwner.avatar ? (
          <Avatar src={siteOwner.avatar} size={48} />
        ) : (
          <Avatar size={48}>{siteOwner.name.slice(0, 1)}</Avatar>
        )}
      </button>

      {drawerMounted && (
        <div className="fixed inset-0 z-[1000] md:hidden">
          <button
            type="button"
            aria-label="关闭站长信息"
            onClick={closeOwnerPanel}
            className={`absolute inset-0 bg-black/25 backdrop-blur-sm transition-opacity duration-300 ease-in ${
              open ? "opacity-100" : "opacity-0"
            }`}
          />
          <aside
            className={`absolute right-0 top-0 flex h-full w-1/2 flex-col overflow-y-auto bg-white px-3 py-5 shadow-2xl transition-transform duration-300 ease-in dark:bg-[#141414] dark:text-gray-100 ${
              open ? "translate-x-0" : "translate-x-full"
            }`}
          >
            <div className="mb-5 flex items-start justify-between gap-2">
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  {siteOwner.avatar ? (
                    <Avatar src={siteOwner.avatar} size={44} />
                  ) : (
                    <Avatar size={44}>{siteOwner.name.slice(0, 1)}</Avatar>
                  )}
                  <div className="min-w-0">
                    <div className="truncate text-sm font-semibold">
                      {siteOwner.name}
                    </div>
                    <div className="mt-0.5 text-xs text-gray-400">站长信息</div>
                  </div>
                </div>
              </div>
              <Button
                type="text"
                shape="circle"
                icon={<CloseOutlined />}
                onClick={closeOwnerPanel}
              />
            </div>

            <div className="rounded-lg border border-gray-100 bg-gray-50 p-3 text-xs leading-5 text-gray-600 dark:border-gray-700 dark:bg-[#232426] dark:text-gray-300">
              {siteOwner.hasError
                ? "网站数据暂时异常"
                : siteOwner.bio || "这个站长还没有填写简介。"}
            </div>

            {statsItems.length > 0 && (
              <div className="mt-4 grid grid-cols-2 gap-2 text-center text-xs">
                {statsItems.map((item) => (
                  <div
                    key={item.label}
                    className="rounded-lg border border-gray-100 bg-gray-50 px-2 py-2 dark:border-gray-700 dark:bg-[#232426]"
                  >
                    <div className="text-gray-400">{item.label}</div>
                    <div className="mt-1 truncate font-semibold text-gray-700 dark:text-gray-200">
                      {item.value}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </aside>
        </div>
      )}
    </>
  );
}
