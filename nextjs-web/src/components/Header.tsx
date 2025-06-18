"use client";
import React from "react";
import { Avatar, Button, Space } from "antd";
import { BulbOutlined, SearchOutlined } from "@ant-design/icons";
import Navbar from "./Navbar";
import Link from "next/link";

const Header: React.FC = () => {
  return (
    <header className="w-full bg-white border-b border-gray-200 px-4 py-2 flex items-center justify-between">
      <div className="flex items-center space-x-4 w-full">
        <Link href="/">
          <Avatar src="/logo.png" alt="logo" size={40} />
        </Link>
        <div className="flex-1 min-w-0">
          <Navbar />
        </div>
      </div>
      <Space size="middle">
        <Button type="text" shape="circle" icon={<BulbOutlined />} aria-label="切换暗黑模式" />
        <Button type="text" shape="circle" icon={<SearchOutlined />} aria-label="搜索" />
      </Space>
    </header>
  );
};

export default Header;
