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
import { getIndexConfig } from "@/src/api/config";
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
  const config = await getIndexConfig();
  const seoConfig = config.seo_meta_config;
  
  return {
    title: seoConfig.title || config.website_config.website_name,
    description: seoConfig.description,
    keywords: seoConfig.keywords,
    authors: [{ name: seoConfig.author || config.website_config.website_owner }],
    robots: seoConfig.robots || "index, follow",
    openGraph: {
      title: seoConfig.og_title || seoConfig.title,
      description: seoConfig.description,
      url: process.env.BASE_HOST,
      images: seoConfig.og_image ? [{ url: process.env.SERVER_HOST + seoConfig.og_image }] : [],
      siteName: config.website_config.website_name,
      type: "website",
    },
    verification: {
      // 百度验证通过 other 字段处理
      google: undefined, // 可以根据需要添加谷歌验证
    },
    other: {
      // 百度站点验证
      ...(seoConfig.baidu_site_verification && {
        'baidu-site-verification': seoConfig.baidu_site_verification,
      }),
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
  const config = await getIndexConfig();
  let stats = undefined;
  try {
    stats = await getWebsiteStats();
  } catch {}
  const configWithStats = { ...config, stats };

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AntdRegistry>
          <AntdThemeProvider>
            <div className="min-h-screen flex flex-col">
              <Header menus={menus} />
              <ConfigToZustand config={configWithStats} />
              <SeoHead config={configWithStats} />
              <main>{children}</main>
              <LogVisitClient />
              <Footer websiteRecords={configWithStats.website_config.website_records || []} />
            </div>
          </AntdThemeProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}
