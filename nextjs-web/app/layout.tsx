import type { Metadata } from "next";
import "./globals.css";
import { ConfigProvider } from "./context/ConfigContext";
import BlogLayout from "./components/BlogLayout";
import StyledComponentsRegistry from "./components/AntdRegistry";
import { AntdRegistry } from '@ant-design/nextjs-registry';

export async function generateMetadata(): Promise<Metadata> {
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || 'http://localhost:8080';
  const url = `${apiHost}/configs/index`;

  try {
    console.log('nihao');
    
    const res = await fetch(url, { cache: 'no-store' });
    if (res.ok) {
      const text = await res.text();
      try {
        const data = JSON.parse(text);
        const websiteInfo = data?.data || {};
        const seo = websiteInfo?.seo_meta_config || {};
        const website = websiteInfo?.website_config || {};

        return {
          title: seo.title || website.website_name || 'fnote',
          description: seo.description || '',
          keywords: seo.keywords || '',
          authors: seo.author ? [{ name: seo.author }] : [],
          robots: seo.robots || 'index, follow',
          icons: {
            icon: website.website_icon?.startsWith('http')
              ? website.website_icon
              : apiHost + (website.website_icon || '/favicon.ico'),
          },
          openGraph: {
            title: seo.og_title || seo.title || website.website_name,
            description: seo.description || '',
            images: seo.og_image ? [seo.og_image] : [],
            url: apiHost,
            siteName: website.website_name || '',
            type: 'website',
          },
          twitter: {
            card: seo.twitter_card || 'summary_large_image',
            title: seo.og_title || seo.title || website.website_name,
            description: seo.description || '',
            images: seo.og_image ? [seo.og_image] : [],
          },
        };
      } catch (jsonErr) {
        console.error('generateMetadata JSON parse error:', jsonErr, text);
      }
    } else {
      console.error(`generateMetadata fetch error: HTTP ${res.status}`);
    }
  } catch (err) {
    console.error('generateMetadata fetch error:', err);
  }
  // 兜底返回空对象
  return {};
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN" suppressHydrationWarning>
      <body className="antialiased">
        <AntdRegistry>
          <ConfigProvider>
            <BlogLayout>{children}</BlogLayout>
          </ConfigProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}
