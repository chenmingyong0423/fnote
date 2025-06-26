"use client";
import React, { useEffect, useState } from "react";
import { ConfigProvider, theme as antdTheme } from "antd";

function getInitialDark() {
  if (typeof window !== "undefined") {
    const saved = localStorage.getItem("theme-dark");
    if (saved === "1") return true;
    if (saved === "0") return false;
    if (document.documentElement.classList.contains("dark")) return true;
    return window.matchMedia("(prefers-color-scheme: dark)").matches;
  }
  return undefined;
}

export function AntdThemeProvider({ children }: { children: React.ReactNode }) {
  const [isDark, setIsDark] = useState<boolean | undefined>(undefined);

  useEffect(() => {
    const checkDark = () => {
      setIsDark(document.documentElement.classList.contains("dark"));
    };
    checkDark();
    const observer = new MutationObserver(checkDark);
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ["class"] });
    return () => observer.disconnect();
  }, []);

  if (isDark === undefined) return null; // 避免初始 SSR/CSR 不一致

  return (
    <ConfigProvider
      theme={{
        algorithm: isDark ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
      }}
    >
      {children}
    </ConfigProvider>
  );
}
