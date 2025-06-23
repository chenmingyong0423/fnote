"use client";
import React from "react";
import { Menu, Spin } from "antd";
import Link from "next/link";
import type { MenuVO } from "../api/category";

const staticNav = [
  { label: "首页", href: "/" },
  { label: "文章导航", href: "/navigation" },
  { label: "关于", href: "/about" },
];

const Navbar: React.FC<{ menus: MenuVO[]; loading?: boolean }> = ({ menus, loading }) => {
  // 组装 items 数组
  const items = [
    {
      key: staticNav[0].href,
      label: <Link href={staticNav[0].href}>{staticNav[0].label}</Link>,
    },
    {
      key: staticNav[1].href,
      label: <Link href={staticNav[1].href}>{staticNav[1].label}</Link>,
    },
    ...(loading
      ? [
          {
            key: "loading",
            label: <Spin size="small" />, // 禁用项
            disabled: true,
          },
        ]
      : menus.map((item) => ({
          key: item.route,
          label: <Link href={`/categories${item.route}`}>{item.name}</Link>,
        }))) ,
    {
      key: staticNav[2].href,
      label: <Link href={staticNav[2].href}>{staticNav[2].label}</Link>,
    },
  ];

  return (
    <Menu
      mode="horizontal"
      selectable={false}
      className="bg-transparent border-none shadow-none"
      items={items}
      style={{ flex: 1, minWidth: 0 }}
    />
  );
};

export default Navbar;
