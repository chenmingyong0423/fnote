"use client";
import { Button, Space } from "antd";
import { BulbOutlined, SearchOutlined } from "@ant-design/icons";
import Link from "next/link";
import { useThemeStore } from "../store/theme";

const HeaderActions = () => {
  const isDark = useThemeStore((s) => s.isDark);
  const toggleDarkMode = useThemeStore((s) => s.toggleDark);

  return (
    <Space size={4} className="md:[&_.ant-space-item:not(:last-child)]:mr-2">
      <Button
        type="text"
        shape="circle"
        icon={<BulbOutlined />}
        aria-label="切换暗黑模式"
        onClick={toggleDarkMode}
        style={{ color: isDark ? "#fadb14" : undefined }}
      />
      <Link href="/search">
        <Button type="text" shape="circle" icon={<SearchOutlined />} aria-label="搜索" />
      </Link>
    </Space>
  );
};

export default HeaderActions;
