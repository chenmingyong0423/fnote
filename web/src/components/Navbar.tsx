"use client";
import React from "react";
import { Menu, Spin } from "antd";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import type { MenuVO } from "../api/category";

const staticNav = [
  { label: "首页", href: "/" },
  { label: "文章导航", href: "/navigation" },
  { label: "友链", href: "/friend" },
  { label: "关于", href: "/about" },
];

const Navbar: React.FC<{ menus: MenuVO[]; loading?: boolean }> = ({ menus, loading }) => {
  const pathname = usePathname();
  const router = useRouter();
  const getSelectedKeys = (pathname: string) => {
    pathname = pathname.replace(/\/$/, ""); // 去除末尾 /
    if (pathname === "" || pathname === "/") return ["/"];
    if (pathname === "/navigation" || pathname === "/friend" || pathname === "/about") {
      return [pathname];
    }
    // 其它如 /categories/frontend => ["categories", "frontend"]
    return pathname.replace(/^\//, "").split("/");
  };
  const [selectedKeys, setSelectedKeys] = React.useState(getSelectedKeys(pathname ?? ""));
  React.useEffect(() => {
    setSelectedKeys(getSelectedKeys(pathname ?? ""));
  }, [pathname]);
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
          label: <Link href={`/categories/${item.route}`}>{item.name}</Link>,
        }))) ,
    {
      key: staticNav[2].href,
      label: <Link href={staticNav[2].href}>{staticNav[2].label}</Link>,
    },
    {
      key: staticNav[3].href,
      label: <Link href={staticNav[3].href}>{staticNav[3].label}</Link>,
    },
  ];

  return (
    <div className="w-full min-w-0 overflow-x-auto overflow-y-hidden [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
    <Menu
      mode="horizontal"
      selectable={false}
      className="min-w-max bg-transparent border-none shadow-none [&_.ant-menu-item]:px-2 md:[&_.ant-menu-item]:px-4"
      items={items}
      style={{ flex: 1, minWidth: 0 }}
      selectedKeys={selectedKeys}
      onSelect={({ keyPath }) => {
        setSelectedKeys(keyPath);
        // 还原为路径
        const url = keyPath[0] === "/" ? "/" : `/${keyPath.reverse().join("/")}`;
        router.push(url);
      }}
    />
    </div>
  );
};

export default Navbar;
