import React from "react";
import { Avatar } from "antd";
import Navbar from "./Navbar";
import Link from "next/link";
import { getMenus } from "../api/category";
import HeaderActions from "./HeaderActions";
import type { WebsiteMetaVO } from "@/src/api/config";

type HeaderProps = {
  websiteMetaConfig?: WebsiteMetaVO;
};

const Header = async ({ websiteMetaConfig }: HeaderProps) => {
  const menus = await getMenus().catch(() => []);

  return (
    <header
      className="w-full bg-white dark:bg-[#141414] border-b border-gray-100 dark:border-[#303030] shadow-sm rounded-xl py-0 px-4 mx-auto max-w-7xl grid grid-cols-12 items-center h-[60px] mt-4 mb-8"
    >
      {/* 左侧 Logo 区 1/12 */}
      <div className="col-span-1 flex items-center">
        <Link href="/">
          <Avatar src={websiteMetaConfig?.website_icon || "/logo.png"} alt="logo" size={40} />
        </Link>
      </div>
      {/* 菜单区 7/12，左对齐，紧挨logo */}
      <div className="col-span-7 flex items-center justify-start min-w-0">
        <Navbar menus={menus} />
      </div>
      {/* 右侧按钮区 4/12 */}
      <div className="col-span-4 flex items-center justify-end gap-2">
        <HeaderActions />
      </div>
    </header>
  );
};

export default Header;
