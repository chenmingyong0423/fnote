import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import React from "react";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import Footer from "../src/components/Footer";
import Header from "../src/components/Header";
import SeoHead from "../src/components/SeoHead";
import { getCommonConfig } from "@/src/api/config";
import { AntdThemeProvider } from "@/src/components/AntdThemeProvider";
import LogVisitClient from "../src/components/LogVisitClient";
import { checkInitialization } from "@/src/api/checkInitialization";
import { redirect } from "next/navigation";
export const dynamic = 'force-dynamic'

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

type InitCache = {
  promise: Promise<void> | null;
  expiresAt: number;
};

const INIT_TTL_MS = 30_000; // 30s，可按需调整
const initCache: InitCache = {
  promise: null,
  expiresAt: 0,
};

async function guardInitialization() {
  const initRes = await checkInitialization();
  if (initRes.code === 0 && initRes.data && !initRes.data.initStatus) {
    const adminHost = process.env.NEXT_PUBLIC_ADMIN_HOST || process.env.ADMIN_HOST;
    if (adminHost) {
      redirect(adminHost);
    }
    throw new Error("Site not initialized and admin host is not configured");
  }
}

function ensureInitialized() {
  const now = Date.now();

  // 缓存未过期，直接复用
  if (initCache.promise && now < initCache.expiresAt) {
    return initCache.promise;
  }

  // 重新发起检查，并刷新过期时间
  initCache.expiresAt = now + INIT_TTL_MS;
  initCache.promise = guardInitialization().catch((err) => {
    // 失败时清空，避免后续一直复用失败 Promise
    initCache.promise = null;
    initCache.expiresAt = 0;
    throw err;
  });

  return initCache.promise;
}

// 动态生成 metadata
export async function generateMetadata(): Promise<Metadata> {
  await ensureInitialized();
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
      images: config.seo_meta.og_image ? [{ url: process.env.NEXT_PUBLIC_SERVER_HOST + config.seo_meta.og_image }] : undefined,
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
  await ensureInitialized();
  const config = await getCommonConfig();

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AntdRegistry>
          <AntdThemeProvider>
            <div className="min-h-screen flex flex-col">
              <Header />
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
