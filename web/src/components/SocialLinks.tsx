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
    "inline-flex h-9 w-9 items-center justify-center rounded-full border border-gray-200 bg-white/60 text-base text-gray-600 transition-colors hover:border-blue-300 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 dark:border-gray-700 dark:bg-white/5 dark:text-gray-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400";

  return (
    <div className="flex flex-wrap items-center justify-center gap-2">
      {items.map((item) => {
        const key = `${item.social_name}-${item.social_value}`;
        const tooltip = item.is_link
          ? item.social_name
          : `${item.social_name}：${item.social_value}`;

        return (
          <Tooltip key={key} title={tooltip}>
            {item.is_link ? (
              <a
                href={resolveHref(item.social_value)}
                target="_blank"
                rel="noopener noreferrer"
                className={itemClass}
                aria-label={`打开${item.social_name}`}
              >
                <SocialIcon cssClass={item.css_class} />
              </a>
            ) : (
              <button
                type="button"
                className={itemClass}
                aria-label={`复制${item.social_name}`}
                onClick={() => handleCopy(item)}
              >
                <SocialIcon cssClass={item.css_class} />
              </button>
            )}
          </Tooltip>
        );
      })}
    </div>
  );
}
