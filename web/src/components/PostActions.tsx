"use client";
import React, { useState, useEffect } from "react";
import { App, Button, Tooltip, Popover } from "antd";
import {
  LikeOutlined,
  MessageOutlined,
  ShareAltOutlined,
  GiftOutlined,
  WechatOutlined,
  LinkOutlined,
  LikeFilled,
} from "@ant-design/icons";
import { QRCodeCanvas } from "qrcode.react";
import { getCommonConfig, type PayInfoConfigVO } from "@/src/api/config";
import { likePost } from "@/src/api/posts";

let payInfoCache: PayInfoConfigVO[] | null = null;
let payInfoRequest: Promise<PayInfoConfigVO[]> | null = null;

function getPayImageSrc(value: string) {
  const src = value.trim();
  if (!src) return src;

  try {
    const url = new URL(src);
    if (url.pathname.startsWith("/static/")) {
      return `${url.pathname}${url.search}${url.hash}`;
    }
  } catch {
    return src;
  }

  return src;
}

async function getCachedPayInfo() {
  if (payInfoCache !== null) return payInfoCache;

  if (!payInfoRequest) {
    payInfoRequest = getCommonConfig()
      .then((config) =>
        (config.pay_info_config ?? []).filter(
          (item) => item.name && item.image,
        ),
      )
      .catch((error) => {
        payInfoRequest = null;
        throw error;
      });
  }

  payInfoCache = await payInfoRequest;
  return payInfoCache;
}

interface PostActionsProps {
  postId: string;
  isLiked?: boolean;
  likeCount?: number;
  commentCount?: number;
}

export const PostActions: React.FC<PostActionsProps> = ({
  postId,
  isLiked = false,
  likeCount = 0,
  commentCount = 0,
}) => {
  const { message } = App.useApp();
  const [liked, setLiked] = useState(isLiked);
  const [currentLikeCount, setCurrentLikeCount] = useState(likeCount);
  const [likeLoading, setLikeLoading] = useState(false);
  const [currentUrl, setCurrentUrl] = useState("");
  const [payInfoList, setPayInfoList] = useState<PayInfoConfigVO[]>(
    payInfoCache ?? [],
  );
  const [payInfoLoading, setPayInfoLoading] = useState(payInfoCache === null);

  useEffect(() => {
    setCurrentUrl(window.location.href);
  }, []);

  useEffect(() => {
    setLiked(isLiked);
    setCurrentLikeCount(likeCount);
  }, [isLiked, likeCount]);

  useEffect(() => {
    let cancelled = false;

    getCachedPayInfo()
      .then((list) => {
        if (!cancelled) {
          setPayInfoList(list);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setPayInfoList([]);
        }
      })
      .finally(() => {
        if (!cancelled) {
          setPayInfoLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  const handleCopy = () => {
    navigator.clipboard.writeText(currentUrl).then(() => {
      message.success("链接已复制");
    });
  };

  const handleLike = async () => {
    if (liked) return;
    setLikeLoading(true);
    try {
      const res = await likePost(postId);
      if (res.code === 0) {
        setLiked(true);
        setCurrentLikeCount((count) => count + 1);
        message.success("点赞成功");
      } else {
        message.error(res.message || "点赞失败");
      }
    } catch (e: unknown) {
      message.error((e instanceof Error ? e.message : String(e)) || "点赞失败");
    } finally {
      setLikeLoading(false);
    }
  };

  return (
    <div className="glass-surface flex flex-row items-center justify-center gap-1.5 md:gap-3 rounded-lg p-2 md:p-3">
      <Tooltip title={liked ? "已点赞" : "点赞"}>
        <Button
          type="text"
          icon={
            liked ? (
              <LikeFilled style={{ color: "#eb2f96" }} />
            ) : (
              <LikeOutlined />
            )
          }
          loading={likeLoading}
          onClick={handleLike}
          disabled={liked}
        >
          {currentLikeCount}
        </Button>
      </Tooltip>
      <Tooltip title="评论">
        <Button type="text" icon={<MessageOutlined />} href="#comments">
          {commentCount}
        </Button>
      </Tooltip>
      <Popover
        placement="bottom"
        content={
          <div className="flex flex-row items-center gap-3">
            <Popover
              placement="right"
              content={<QRCodeCanvas value={currentUrl} size={120} />}
              trigger="hover"
            >
              <Tooltip title="微信">
                <Button type="text" shape="circle" icon={<WechatOutlined />} />
              </Tooltip>
            </Popover>
            <Tooltip title="复制链接">
              <Button
                type="text"
                shape="circle"
                icon={<LinkOutlined />}
                onClick={handleCopy}
              />
            </Tooltip>
          </div>
        }
        trigger="hover"
      >
        <Tooltip title="分享">
          <Button type="text" icon={<ShareAltOutlined />} />
        </Tooltip>
      </Popover>
      <Popover
        placement="bottom"
        content={
          <div className="flex min-w-24 flex-col gap-3">
            {payInfoLoading ? (
              <span className="text-xs text-gray-500 dark:text-gray-400">
                加载中...
              </span>
            ) : payInfoList.length > 0 ? (
              payInfoList.map((item) => (
                <div
                  key={`${item.name}-${item.image}`}
                  className="flex flex-col items-center"
                >
                  <img
                    src={getPayImageSrc(item.image)}
                    alt={item.name}
                    width={96}
                    height={96}
                    loading="lazy"
                    decoding="async"
                    className="h-24 w-24 rounded border border-gray-200 object-contain dark:border-gray-700"
                  />
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    {item.name}
                  </span>
                </div>
              ))
            ) : (
              <span className="text-xs text-gray-500 dark:text-gray-400">
                暂无赞赏码
              </span>
            )}
          </div>
        }
        trigger="hover"
      >
        <Tooltip title="赞赏">
          <Button type="text" icon={<GiftOutlined />} />
        </Tooltip>
      </Popover>
    </div>
  );
};
