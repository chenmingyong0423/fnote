import React from "react";
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
      className="w-[calc(100%-2rem)] md:w-full bg-white/95 dark:bg-[#141414]/95 border border-gray-200/80 dark:border-[#303030] shadow-sm shadow-gray-200/70 dark:shadow-none rounded-lg px-3 md:px-4 mx-auto max-w-7xl grid grid-cols-[auto_1fr_auto] md:grid-cols-12 items-center gap-2 md:gap-0 min-h-[56px] md:h-[60px] mt-3 md:mt-4 mb-5 md:mb-8 backdrop-blur"
    >
      {/* 左侧 Logo 区 1/12 */}
      <div className="md:col-span-1 flex items-center">
        <Link href="/">
          <img
            src={websiteMetaConfig?.website_icon || "/logo.png"}
            alt="logo"
            width={36}
            height={36}
            loading="lazy"
            className="h-9 w-9 rounded-full object-cover ring-2 ring-gray-100 transition-transform hover:scale-105 dark:ring-[#303030] md:h-10 md:w-10"
          />
        </Link>
      </div>
      {/* 菜单区 7/12，左对齐，紧挨logo */}
      <div className="min-w-0 overflow-visible md:col-span-7 flex items-center justify-start">
        <Navbar menus={menus} />
      </div>
      {/* 右侧按钮区 4/12 */}
      <div className="md:col-span-4 flex items-center justify-end gap-1 md:gap-2">
        <HeaderActions />
      </div>
    </header>
  );
};

export default Header;
