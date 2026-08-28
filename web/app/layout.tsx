import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "antd/dist/reset.css";
import "./globals.css";
import React from "react";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import Footer from "../src/components/Footer";
import Header from "../src/components/Header";
import { DEFAULT_COMMON_CONFIG, getCommonConfig } from "@/src/api/config";
import { AntdThemeProvider } from "@/src/components/AntdThemeProvider";
import LogVisitClient from "../src/components/LogVisitClient";
import { checkInitialization } from "@/src/api/checkInitialization";
import { redirect } from "next/navigation";
import { isBackendUnavailableError } from "@/src/utils/http";
import { resolvePublicUrl } from "@/src/utils/publicUrl";
import BackTopButton from "@/src/components/BackTopButton";
import InteractiveBackdrop from "@/src/components/InteractiveBackdrop";
import { getSiteUrl, normalizeRobots } from "@/src/utils/seo";
import { headers } from "next/headers";
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
  const requestHeaders = await headers();
  const pathname = requestHeaders.get("x-fnote-pathname") || "/";
  const siteUrl = getSiteUrl();

  try {
    await ensureInitialized();
    const config = await getCommonConfig();
    return {
      metadataBase: siteUrl,
      title: config.seo_meta.title || config.website_meta.website_name,
      description: config.seo_meta.description,
      keywords: config.seo_meta.keywords,
      authors: [{ name: config.seo_meta.author || config.website_meta.website_owner }],
      robots: normalizeRobots(config.seo_meta.robots),
      alternates: {
        canonical: pathname,
      },
      icons: config.website_meta.website_icon
        ? { icon: resolvePublicUrl(config.website_meta.website_icon) }
        : undefined,
      openGraph: {
        title: config.seo_meta.og_title || config.website_meta.website_name,
        description: config.seo_meta.description,
        url: pathname,
        images: config.seo_meta.og_image ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }] : undefined,
        siteName: config.website_meta.website_name,
        locale: "zh_CN",
        type: "website",
      },
      twitter: {
        card: "summary_large_image",
        title: config.seo_meta.og_title || config.website_meta.website_name,
        description: config.seo_meta.description,
        images: config.seo_meta.og_image
          ? [resolvePublicUrl(config.seo_meta.og_image)!]
          : undefined,
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
  } catch (error) {
    if (isBackendUnavailableError(error)) {
      return {
        metadataBase: siteUrl,
        title: DEFAULT_COMMON_CONFIG.seo_meta.title,
        description: DEFAULT_COMMON_CONFIG.seo_meta.description,
        alternates: {
          canonical: pathname,
        },
      };
    }
    throw error;
  }
}

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  let config = DEFAULT_COMMON_CONFIG;
  let hasSiteIssue = false;
  try {
    await ensureInitialized();
    config = await getCommonConfig();
  } catch (error) {
    if (!isBackendUnavailableError(error)) {
      throw error;
    }
    hasSiteIssue = true;
  }

  return (
    <html lang="zh-CN" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var saved=localStorage.getItem("theme-dark");var prefersDark=window.matchMedia&&window.matchMedia("(prefers-color-scheme: dark)").matches;var dark=saved==="1"||(saved===null&&prefersDark);document.documentElement.classList.toggle("dark",dark);}catch(e){}})();`,
          }}
        />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AntdRegistry>
          <AntdThemeProvider>
            <InteractiveBackdrop />
            <div className="relative z-10 min-h-screen flex flex-col">
              <Header websiteMetaConfig={config.website_meta} />
              {hasSiteIssue && (
                <div className="mx-auto mb-5 md:mb-6 w-[calc(100%-2rem)] md:w-full max-w-7xl rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
                  网站数据暂时异常，部分内容可能无法显示，请稍后再试。
                </div>
              )}
              <main>{children}</main>
              <LogVisitClient />
              <Footer websiteRecords={config.records || []} />
              <BackTopButton />
            </div>
          </AntdThemeProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}
