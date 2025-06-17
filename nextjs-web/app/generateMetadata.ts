import { Metadata } from "next";

export async function generateMetadata(): Promise<Metadata> {
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || 'http://localhost:8080';
  const url = `${apiHost}/api/configs/index`;
  let websiteInfo: any = {};
  let icon = '/favicon.ico';
  let title = 'fnote';
  let description = 'fnote';

  try {
    const res = await fetch(url, { cache: 'no-store' });
    const text = await res.text();
    // 尝试解析 JSON
    try {
      const data = JSON.parse(text);
      websiteInfo = data?.data?.website_config || {};
      icon = (websiteInfo as any).website_icon || '/favicon.ico';
      title = (websiteInfo as any).website_name || 'fnote';
      description = (websiteInfo as any).website_desc || 'fnote';
    } catch (jsonErr) {
      // 不是 JSON，打印内容方便排查
      console.error('generateMetadata JSON parse error:', jsonErr, text);
    }
  } catch (err) {
    // fetch 失败
    console.error('generateMetadata fetch error:', err);
  }

  return {
    title,
    description,
    icons: {
      icon: icon.startsWith('http') ? icon : apiHost + icon,
    },
  };
}
