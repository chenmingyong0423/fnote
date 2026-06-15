"use client";

import { MoreOutlined } from "@ant-design/icons";
import { Button, Dropdown, Menu, Spin } from "antd";
import type { MenuProps } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react";
import type { MenuVO } from "../api/category";

const staticNav = [
  { label: "首页", href: "/" },
  { label: "文章导航", href: "/navigation" },
  { label: "友链", href: "/friend" },
  { label: "关于", href: "/about" },
];

type NavEntry = {
  disabled?: boolean;
  href: string;
  key: string;
  label: React.ReactNode;
};
type MenuItems = NonNullable<MenuProps["items"]>;

const Navbar: React.FC<{ menus: MenuVO[]; loading?: boolean }> = ({
  menus,
  loading,
}) => {
  const pathname = usePathname();
  const [mobileMoreOpen, setMobileMoreOpen] = React.useState(false);
  const [mobileVisibleCount, setMobileVisibleCount] = React.useState(2);
  const mobileContainerRef = React.useRef<HTMLDivElement | null>(null);
  const mobileMeasureRef = React.useRef<HTMLDivElement | null>(null);

  const getSelectedKeys = (currentPathname: string) => {
    const normalizedPathname = currentPathname.replace(/\/$/, "");
    if (normalizedPathname === "" || normalizedPathname === "/") return ["/"];
    if (
      normalizedPathname === "/navigation" ||
      normalizedPathname === "/friend" ||
      normalizedPathname === "/about"
    ) {
      return [normalizedPathname];
    }
    return normalizedPathname.replace(/^\//, "").split("/");
  };

  const [selectedKeys, setSelectedKeys] = React.useState(
    getSelectedKeys(pathname ?? "")
  );

  React.useEffect(() => {
    setSelectedKeys(getSelectedKeys(pathname ?? ""));
    setMobileMoreOpen(false);
  }, [pathname]);

  const navEntries: NavEntry[] = [
    {
      key: staticNav[0].href,
      href: staticNav[0].href,
      label: staticNav[0].label,
    },
    {
      key: staticNav[1].href,
      href: staticNav[1].href,
      label: staticNav[1].label,
    },
    ...(loading
      ? [
          {
            key: "loading",
            href: "",
            label: <Spin size="small" />,
            disabled: true,
          },
        ]
      : menus.map((item) => ({
          key: item.route,
          href: `/categories/${item.route}`,
          label: item.name,
        }))),
    {
      key: staticNav[2].href,
      href: staticNav[2].href,
      label: staticNav[2].label,
    },
    {
      key: staticNav[3].href,
      href: staticNav[3].href,
      label: staticNav[3].label,
    },
  ];
  const navSignature = navEntries.map((item) => item.key).join("|");

  const hrefByKey = new Map(navEntries.map((item) => [item.key, item.href]));
  const toMenuItems = (entries: NavEntry[], useLink = true): MenuItems =>
    entries.map((item) => ({
      key: item.key,
      disabled: item.disabled,
      label:
        useLink && item.href ? (
          <Link href={item.href} target="_blank" rel="noopener noreferrer">
            {item.label}
          </Link>
        ) : (
          item.label
        ),
    }));

  const menuItems = toMenuItems(navEntries);
  const mobilePrimaryEntries = navEntries.slice(0, mobileVisibleCount);
  const mobileMoreItems = toMenuItems(
    navEntries.slice(mobileVisibleCount),
    false
  );

  React.useEffect(() => {
    const container = mobileContainerRef.current;
    const measure = mobileMeasureRef.current;
    if (!container || !measure) return;

    const calculateVisibleCount = () => {
      const containerWidth = container.clientWidth;
      const itemWidths = Array.from(
        measure.querySelectorAll<HTMLElement>("[data-nav-measure='item']")
      ).map((element) => Math.ceil(element.getBoundingClientRect().width));
      const moreWidth = Math.ceil(
        measure
          .querySelector<HTMLElement>("[data-nav-measure='more']")
          ?.getBoundingClientRect().width || 40
      );

      if (containerWidth <= 0 || itemWidths.length === 0) return;

      let usedWidth = 0;
      let nextVisibleCount = 0;

      for (let index = 0; index < itemWidths.length; index += 1) {
        const hasRemainingItems = index < itemWidths.length - 1;
        const reservedMoreWidth = hasRemainingItems ? moreWidth : 0;
        if (
          usedWidth + itemWidths[index] + reservedMoreWidth >
          containerWidth
        ) {
          break;
        }

        usedWidth += itemWidths[index];
        nextVisibleCount = index + 1;
      }

      nextVisibleCount = Math.max(1, nextVisibleCount);
      setMobileVisibleCount((current) =>
        current === nextVisibleCount ? current : nextVisibleCount
      );
    };

    calculateVisibleCount();

    const resizeObserver = new ResizeObserver(calculateVisibleCount);
    resizeObserver.observe(container);

    return () => resizeObserver.disconnect();
  }, [navSignature]);

  const openByKey = (key: string) => {
    const url = hrefByKey.get(key);
    if (url) {
      window.open(url, "_blank", "noopener,noreferrer");
    }
  };

  const handleSelect = ({ keyPath }: { keyPath: string[] }) => {
    setSelectedKeys(keyPath);
  };

  return (
    <div className="w-full min-w-0">
      <div
        ref={mobileContainerRef}
        className="relative flex w-full min-w-0 items-center md:hidden"
      >
        <nav className="flex min-w-0 flex-1 flex-nowrap items-center overflow-hidden">
          {mobilePrimaryEntries.map((item) =>
            item.href ? (
              <Link
                key={item.key}
                href={item.href}
                target="_blank"
                rel="noopener noreferrer"
                className={`flex h-[46px] shrink-0 items-center whitespace-nowrap px-2 text-sm transition-colors ${
                  selectedKeys.includes(item.key)
                    ? "text-blue-600 dark:text-blue-400"
                    : "text-gray-700 dark:text-gray-200"
                }`}
              >
                {item.label}
              </Link>
            ) : (
              <span
                key={item.key}
                className="flex h-[46px] shrink-0 items-center whitespace-nowrap px-2 text-sm"
              >
                {item.label}
              </span>
            )
          )}
        </nav>
        {mobileMoreItems.length > 0 && (
          <Dropdown
            trigger={["click"]}
            open={mobileMoreOpen}
            onOpenChange={setMobileMoreOpen}
            placement="bottomRight"
            menu={{
              items: mobileMoreItems,
              selectable: false,
              onClick: ({ key }) => {
                setMobileMoreOpen(false);
                openByKey(String(key));
              },
            }}
          >
            <Button
              type="text"
              shape="circle"
              aria-label="更多导航"
              icon={<MoreOutlined />}
              onClick={(event) => event.preventDefault()}
            />
          </Dropdown>
        )}
        <div
          ref={mobileMeasureRef}
          aria-hidden="true"
          className="pointer-events-none absolute -left-[9999px] top-0 flex h-0 overflow-hidden opacity-0"
        >
          {navEntries.map((item) => (
            <span
              key={item.key}
              data-nav-measure="item"
              className="whitespace-nowrap px-2 text-sm leading-[46px]"
            >
              {item.label}
            </span>
          ))}
          <span
            data-nav-measure="more"
            className="inline-flex h-8 w-10 shrink-0"
          />
        </div>
      </div>

      <div className="hidden w-full min-w-0 md:block">
        <Menu
          mode="horizontal"
          selectable={false}
          triggerSubMenuAction="click"
          className="w-full bg-transparent border-none shadow-none md:[&_.ant-menu-item]:px-4"
          items={menuItems}
          style={{ flex: 1, minWidth: 0 }}
          selectedKeys={selectedKeys}
          onSelect={({ keyPath }) => handleSelect({ keyPath })}
        />
      </div>
    </div>
  );
};

export default Navbar;
