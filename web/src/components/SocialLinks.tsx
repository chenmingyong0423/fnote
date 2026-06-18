"use client";

import { App, Tooltip } from "antd";
import {
  BilibiliOutlined,
  FacebookOutlined,
  GithubOutlined,
  InstagramOutlined,
  LinkOutlined,
  QqOutlined,
  WechatOutlined,
  XOutlined,
  YoutubeOutlined,
  ZhihuOutlined,
} from "@ant-design/icons";
import type { SocialInfoVO } from "../api/config";

function SocialIcon({ cssClass }: { cssClass: string }) {
  if (cssClass.includes("x-twitter")) return <XOutlined />;
  if (cssClass.includes("facebook")) return <FacebookOutlined />;
  if (cssClass.includes("instagram")) return <InstagramOutlined />;
  if (cssClass.includes("youtube")) return <YoutubeOutlined />;
  if (cssClass.includes("bilibili")) return <BilibiliOutlined />;
  if (cssClass.includes("qq")) return <QqOutlined />;
  if (cssClass.includes("github") || cssClass.includes("square-git")) {
    return <GithubOutlined />;
  }
  if (cssClass.includes("weixin")) return <WechatOutlined />;
  if (cssClass.includes("zhihu")) return <ZhihuOutlined />;
  return <LinkOutlined />;
}

function resolveHref(value: string) {
  if (/^(https?:\/\/|mailto:|tel:)/i.test(value)) return value;
  return `https://${value}`;
}

async function copyText(value: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(value);
    return;
  }

  const textarea = document.createElement("textarea");
  textarea.value = value;
  textarea.style.position = "fixed";
  textarea.style.opacity = "0";
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand("copy");
  textarea.remove();
}

export default function SocialLinks({ items }: { items?: SocialInfoVO[] }) {
  const { message } = App.useApp();
  if (!items?.length) return null;

  const handleCopy = async (item: SocialInfoVO) => {
    try {
      await copyText(item.social_value);
      message.success(`${item.social_name}已复制`);
    } catch {
      message.error("复制失败，请手动复制");
    }
  };

  const itemClass =
    "group relative inline-flex h-9 min-w-9 max-w-40 items-center overflow-hidden rounded-full border border-gray-200 bg-white/95 px-[9px] text-base !text-gray-600 shadow-sm transition-colors hover:border-blue-300 hover:bg-blue-50 hover:!text-blue-600 focus-visible:border-blue-300 focus-visible:bg-blue-50 focus-visible:!text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:!text-gray-200 dark:shadow-black/20 dark:hover:border-blue-400 dark:hover:bg-gray-700 dark:hover:!text-blue-300 dark:focus-visible:border-blue-400 dark:focus-visible:bg-gray-700 dark:focus-visible:!text-blue-300 dark:focus-visible:ring-blue-400";
  const labelClass =
    "ml-0 max-w-0 overflow-hidden whitespace-nowrap text-xs font-medium opacity-0 transition-[max-width,margin,opacity] duration-[360ms] ease-[cubic-bezier(0.22,1,0.36,1)] group-hover:ml-2 group-hover:max-w-28 group-hover:opacity-100 group-focus-visible:ml-2 group-focus-visible:max-w-28 group-focus-visible:opacity-100";

  return (
    <div className="flex flex-wrap items-center justify-center gap-2">
      {items.map((item) => {
        const key = `${item.social_name}-${item.social_value}`;

        return (
          <span key={key} className="inline-flex h-9 shrink-0">
            {item.is_link ? (
              <a
                href={resolveHref(item.social_value)}
                target="_blank"
                rel="noopener noreferrer"
                className={itemClass}
                aria-label={`打开${item.social_name}`}
              >
                <span className="inline-flex shrink-0">
                  <SocialIcon cssClass={item.css_class} />
                </span>
                <span className={labelClass}>{item.social_name}</span>
              </a>
            ) : (
              <Tooltip
                title={`${item.social_name}：${item.social_value}`}
                placement="top"
              >
                <button
                  type="button"
                  className={itemClass}
                  aria-label={`复制${item.social_name}`}
                  onClick={() => handleCopy(item)}
                >
                  <span className="inline-flex shrink-0">
                    <SocialIcon cssClass={item.css_class} />
                  </span>
                  <span className={labelClass}>{item.social_name}</span>
                </button>
              </Tooltip>
            )}
          </span>
        );
      })}
    </div>
  );
}
