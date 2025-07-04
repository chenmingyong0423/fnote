import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import React from "react";
import { AntdRegistry } from "@ant-design/nextjs-registry";
// import {App as AntdApp } from "antd";
import Footer from "../src/components/Footer";
import Header from "../src/components/Header";
import ConfigToZustand from "../src/components/ConfigToZustand";
import SeoHead from "../src/components/SeoHead";
import { getCommonConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { getMenus } from "@/src/api/category";
import { AntdThemeProvider } from "@/src/components/AntdThemeProvider";
import LogVisitClient from "../src/components/LogVisitClient";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

// 动态生成 metadata
export async function generateMetadata(): Promise<Metadata> {
  const config = await getCommonConfig();

  return {
    title: config.seo_meta.title || config.website_meta.website_name,
    description: config.seo_meta.description,
    keywords: config.seo_meta.keywords,
    authors: [{ name: config.seo_meta.author || config.website_meta.website_owner }],
    robots: config.seo_meta.robots || "index, follow",
    openGraph: {
      title: config.seo_meta.og_title || config.website_meta.website_name,
      description: config.seo_meta.description,
      url: process.env.BASE_HOST,
      images: config.seo_meta.og_image ? [{ url: process.env.SERVER_HOST + config.seo_meta.og_image }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
    verification: {
      // 百度验证通过 other 字段处理
      google: undefined, // 可以根据需要添加谷歌验证
    },
    other: {
      // 第三方站点验证
      ...config.third_party_site_verification.reduce((acc, item) => {
        acc[item.key] = item.value;
        return acc;
      }, {} as Record<string, string>),
    },
  };
}

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const menus = await getMenus();  // SSR 获取配置信息
  const config = await getCommonConfig();

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AntdRegistry>
          <AntdThemeProvider>
            <div className="min-h-screen flex flex-col">
              <Header menus={menus} />
              <ConfigToZustand config={config} />
              <SeoHead config={config} />
              <main>{children}</main>
              <LogVisitClient />
              <Footer websiteRecords={config.records || []} />
            </div>
          </AntdThemeProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}
