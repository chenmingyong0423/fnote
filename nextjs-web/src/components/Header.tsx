"use client";
import React from "react";
import { Avatar, Button, Space } from "antd";
import { BulbOutlined, SearchOutlined } from "@ant-design/icons";
import Navbar from "./Navbar";
import Link from "next/link";
import { useConfigStore } from "../store/config";
import type { MenuVO } from "../api/category";

const Header: React.FC<{ menus: MenuVO[] }> = ({ menus }) => {
  const websiteConfig = useConfigStore((s) => s.config?.website_config);
  return (
    <header
      className="w-full bg-white border-b border-gray-100 shadow-sm rounded-xl py-0 px-4 mx-auto max-w-7xl grid grid-cols-12 items-center h-[60px] mt-4"
      style={{ minHeight: 60 }}
    >
      {/* 左侧 Logo 区 1/12 */}
      <div className="col-span-1 flex items-center">
        <Link href="/">
          <Avatar src={websiteConfig?.website_icon || "/logo.png"} alt="logo" size={40} />
        </Link>
      </div>
      {/* 菜单区 7/12，左对齐，紧挨logo */}
      <div className="col-span-7 flex items-center justify-start min-w-0">
        <Navbar menus={menus} />
      </div>
      {/* 右侧按钮区 4/12 */}
      <div className="col-span-4 flex items-center justify-end gap-2">
        <Space size="middle">
          <Button type="text" shape="circle" icon={<BulbOutlined />} aria-label="切换暗黑模式" />
          <Button type="text" shape="circle" icon={<SearchOutlined />} aria-label="搜索" />
        </Space>
      </div>
    </header>
  );
};

export default Header;
