"use client";

import { FloatButton } from "antd";

export default function BackTopButton() {
  return (
    <FloatButton.BackTop
      tooltip="返回顶部"
      visibilityHeight={240}
      style={{ insetInlineEnd: 24, insetBlockEnd: 32 }}
    />
  );
}
