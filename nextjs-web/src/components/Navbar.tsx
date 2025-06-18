"use client";
import React, { useEffect, useState } from "react";
import { Menu, Spin } from "antd";
import Link from "next/link";
import { getMenus, MenuVO } from "../api/category";

const staticNav = [
  { label: "首页", href: "/" },
  { label: "关于", href: "/about" },
];

const Navbar: React.FC = () => {
  const [menus, setMenus] = useState<MenuVO[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getMenus()
      .then((data) => setMenus(Array.isArray(data) ? data : []))
      .finally(() => setLoading(false));
  }, []);

  // 组装 items 数组
  const items = [
    {
      key: staticNav[0].href,
      label: <Link href={staticNav[0].href}>{staticNav[0].label}</Link>,
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
      key: staticNav[1].href,
      label: <Link href={staticNav[1].href}>{staticNav[1].label}</Link>,
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
