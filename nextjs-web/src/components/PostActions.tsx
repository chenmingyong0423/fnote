"use client";
import React, { useState } from "react";
import { Button, Tooltip, Popover, message } from "antd";
import { LikeOutlined, MessageOutlined, ShareAltOutlined, GiftOutlined, WechatOutlined, QqOutlined, LinkOutlined, LikeFilled } from "@ant-design/icons";
import { QRCodeCanvas } from "qrcode.react";
import { likePost } from "@/src/api/posts";
import '@ant-design/v5-patch-for-react-19';

const rewardList = [
  { name: "微信赞赏码", img: "/file.svg" },
  { name: "支付宝赞赏码", img: "/globe.svg" },
];

export const PostActions: React.FC<{ postId: string; isLiked?: boolean }> = ({ postId, isLiked = false }) => {
  const [liked, setLiked] = useState(isLiked);
  const [likeLoading, setLikeLoading] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(window.location.href).then(() => {
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
        message.success("点赞成功");
      } else {
        message.error(res.message || "点赞失败");
      }
    } catch (e) {
      message.error("点赞失败");
    } finally {
      setLikeLoading(false);
    }
  };

  return (
      <div className="flex flex-row items-center justify-center gap-3 bg-white dark:bg-[#232426] rounded-lg p-3 shadow-sm border border-gray-100 dark:border-gray-700">
        <Tooltip title={liked ? "已点赞" : "点赞"}>
          <Button
            type="text"
            icon={liked ? <LikeFilled style={{ color: '#eb2f96' }} /> : <LikeOutlined />}
            loading={likeLoading}
            onClick={handleLike}
            disabled={liked}
          />
        </Tooltip>
        <Tooltip title="评论">
          <Button type="text" icon={<MessageOutlined />} href="#comments" />
        </Tooltip>
        <Popover
          placement="bottom"
          content={
            <div className="flex flex-row items-center gap-3">
              <Popover
                placement="right"
                content={<QRCodeCanvas value={window.location.href} size={120} />}
                trigger="hover"
              >
                <Tooltip title="微信">
                  <Button type="text" shape="circle" icon={<WechatOutlined />} />
                </Tooltip>
              </Popover>
              <Tooltip title="复制链接">
                <Button type="text" shape="circle" icon={<LinkOutlined />} onClick={handleCopy} />
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
            <div className="flex flex-col gap-2">
              {rewardList.map(item => (
                <div key={item.name} className="flex flex-col items-center">
                  <img src={item.img} alt={item.name} className="w-20 h-20 object-contain rounded border mb-1" />
                  <span className="text-xs text-gray-500 dark:text-gray-400">{item.name}</span>
                </div>
              ))}
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
