import type { Metadata } from "next";
import "./globals.css";
import { ConfigProvider } from "./context/ConfigContext";
import BlogLayout from "./components/BlogLayout";
import StyledComponentsRegistry from "./components/AntdRegistry";
import { AntdRegistry } from '@ant-design/nextjs-registry';

export const metadata: Metadata = {
  title: "fnote",
  description: "fnote",
};

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
