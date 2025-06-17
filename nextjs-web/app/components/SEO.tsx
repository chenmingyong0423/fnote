'use client';

import { useEffect } from 'react';
import Script from 'next/script';
import { useConfigStore } from '../store/configStore';

interface SEOProps {
  title?: string;
  description?: string;
  keywords?: string;
  ogTitle?: string;
  ogDescription?: string;
  ogImage?: string;
  noindex?: boolean;
}

const SEO: React.FC<SEOProps> = ({
  title,
  description,
  keywords,
  ogTitle,
  ogDescription,
  ogImage,
  noindex
}) => {
  // 直接用 Zustand 获取最新全局配置
  const store = useConfigStore();
  
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || '';
  const domain = process.env.NEXT_PUBLIC_DOMAIN || '';

  // 兼容原 config 结构，显式类型声明，避免类型报错
  const config = {
    website_info: store.website_info?.website_config as {
      website_name?: string;
      website_icon?: string;
    } || {},
    seo_meta_config: store.website_info?.seo_meta_config as {
      title?: string;
      description?: string;
      keywords?: string;
      author?: string;
      robots?: string;
      og_title?: string;
      og_image?: string;
    } || {},
    metaVerificationList: store.website_info?.third_party_site_verification?.map((item: any) => ({
      name: item.key,
      content: item.value
    })) || [],
  };
  // 获取最终的元数据值，优先使用传入的值，其次使用配置中的值
  const finalTitle = title || config.seo_meta_config.title || config.website_info.website_name || 'fnote';
  const finalDescription = description || config.seo_meta_config.description || '';
  
  const finalKeywords = keywords || config.seo_meta_config.keywords || '';
  const finalOgTitle = ogTitle || config.seo_meta_config.og_title || finalTitle;
  const finalOgDescription = ogDescription || config.seo_meta_config.description || '';
  const finalOgImage = ogImage || config.seo_meta_config.og_image || '';
  
  // 构建 JSON-LD 结构化数据
  const jsonLd = JSON.stringify({
    "@context": "https://schema.org",
    "@type": "WebSite",
    "name": finalTitle,
    "url": domain,
  });

  return (
    <>
      <Script
        id="json-ld"
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: jsonLd }}
      />
      {/* 其它 head/meta 相关内容已由 generateMetadata 负责，这里无需再设置 */}
    </>
  );
};

export default SEO;