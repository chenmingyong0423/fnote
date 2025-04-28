'use client';

import Image from 'next/image';
import { useConfigContext } from '../context/ConfigContext';

export default function Profile() {
  const config = useConfigContext();
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || '';
  
  // 获取网站信息
  const websiteInfo = config.website_info || {};
  // 获取统计数据
  const statsInfo = config.website_count_stats || {
    post_count: 0,
    category_count: 0,
    tag_count: 0,
    comment_count: 0,
    like_count: 0,
    view_count: 0
  };
  
  // 图标获取函数（保留原功能，但可能在实际项目中不需要）
  const getIcon = (icon: string): string => {
    switch (icon) {
      case "i-fa6-brands:x-twitter":
        return "i-fa6-brands:x-twitter";
      case "i-fa6-brands:facebook":
        return "i-fa6-brands:facebook";
      case "i-fa6-brands:instagram":
        return "i-fa6-brands:instagram";
      case "i-fa6-brands:youtube":
        return "i-fa6-brands:youtube";
      case "i-fa6-brands:bilibili":
        return "i-fa6-brands:bilibili";
      case "i-fa6-brands:qq":
        return "i-fa6-brands:qq";
      case "i-fa6-brands:github":
        return "i-fa6-brands:github";
      case "i-fa6-brands:square-git":
        return "i-fa6-brands:square-git";
      case "i-fa6-brands:weixin":
        return "i-fa6-brands:weixin";
      case "i-fa6-brands:zhihu":
        return "i-fa6-brands:zhihu";
      case "i-bi:link-45deg":
        return "i-bi:link-45deg";
    }
    return "";
  };

  return (
    <div className="flex flex-col items-center justify-center bg-white p-10 rounded-md dark:text-dtc md:dark:bg-gray-800 ease-linear duration-100 md:shadow-md md:hover:-translate-y-2 sm:p-5 sm:bg-transparent">
      <div className="avatar">
        <img
          src={`${apiHost}${websiteInfo.website_owner_avatar}`}
          alt={websiteInfo.website_owner || "站长头像"}
          className="w-[100px] h-[100px] rounded-full mx-5 cursor-pointer hover:rotate-[360deg] ease-out duration-1000 lg:mr-0"
        />
      </div>
      <div className="introduction flex flex-col items-center justify-center p-5">
        <span className="text-[1.5em] mb-2">{websiteInfo.website_owner}</span>
        <span className="text-gray-500 mb-2">{websiteInfo.website_owner_profile}</span>
      </div>
      <div className="flex items-center justify-between border-t w-full text-gray-500 border-t-[1px] border-t-gray-200 border-solid pt-5 mb-5">
        <div className="flex flex-col items-center justify-center w-[33%]">
          <span className="mb-1">{statsInfo.post_count}</span>
          <span className="">文章</span>
        </div>
        <div className="flex flex-col items-center justify-center w-[33%] border-x-[1px] border-x-gray-200 border-x-solid">
          <span className="mb-1">{statsInfo.category_count}</span>
          <span className="">分类</span>
        </div>
        <div className="flex flex-col items-center justify-center w-[33%]">
          <span className="mb-1">{statsInfo.tag_count}</span>
          <span className="">标签</span>
        </div>
      </div>
      <div className="flex items-center justify-between border-t w-full text-gray-500 pb-5 border-b-[1px] border-b-gray-200 border-b-solid">
        <div className="flex flex-col items-center justify-center w-[33%]">
          <span className="mb-1">{statsInfo.comment_count}</span>
          <span className="">评论</span>
        </div>
        <div className="flex flex-col items-center justify-center w-[33%] border-x-[1px] border-x-gray-200 border-x-solid">
          <span className="mb-1">{statsInfo.like_count}</span>
          <span className="">点赞</span>
        </div>
        <div className="flex flex-col items-center justify-center w-[33%]">
          <span className="mb-1">{statsInfo.website_view_count}</span>
          <span className="">浏览量</span>
        </div>
      </div>
    </div>
  );
}