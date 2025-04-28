'use client';

import { useEffect } from 'react';
import Head from 'next/head';
import Script from 'next/script';
import { useConfigContext } from '../context/ConfigContext';

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
  const config = useConfigContext();
  const apiHost = process.env.NEXT_PUBLIC_API_HOST || '';
  const domain = process.env.NEXT_PUBLIC_DOMAIN || '';

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
    "name": config.website_info.website_name || 'fnote',
    "url": domain,
  });

  // 使用 useEffect 更新文档标题
  useEffect(() => {
    document.title = finalTitle;
  }, [finalTitle]);

  return (
    <>
      <Script
        id="json-ld"
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: jsonLd }}
      />
      <meta name="description" content={finalDescription} />
      <meta name="keywords" content={finalKeywords} />
      <meta name="author" content={config.seo_meta_config.author || 'fnote'} />
      <meta name="robots" content={noindex ? 'noindex, nofollow' : (config.seo_meta_config.robots || 'index, follow')} />
      
      {/* Open Graph / Facebook */}
      <meta property="og:type" content="website" />
      <meta property="og:title" content={finalOgTitle} />
      <meta property="og:description" content={finalOgDescription} />
      <meta property="og:url" content={domain} />
      {finalOgImage && (
        <meta 
          property="og:image" 
          content={finalOgImage.startsWith('http') ? finalOgImage : `${apiHost}${finalOgImage}`} 
        />
      )}
      
      {/* Twitter */}
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={finalOgTitle} />
      <meta name="twitter:description" content={finalOgDescription} />
      {finalOgImage && (
        <meta 
          name="twitter:image" 
          content={finalOgImage.startsWith('http') ? finalOgImage : `${apiHost}${finalOgImage}`} 
        />
      )}
      
      {/* 第三方站点验证 */}
      {config.metaVerificationList?.map((meta, index) => (
        <meta key={index} name={meta.name} content={meta.content} />
      ))}

      {/* 网站图标 */}
      {config.website_info.website_icon && (
        <link 
          rel="icon" 
          type="image/x-icon" 
          href={config.website_info.website_icon.startsWith('http') 
            ? config.website_info.website_icon 
            : `${apiHost}${config.website_info.website_icon}`
          } 
        />
      )}
    </>
  );
};

export default SEO;