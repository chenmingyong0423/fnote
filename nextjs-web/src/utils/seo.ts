import type { Metadata } from 'next';
import type { IndexConfigVO } from '../api/config';

interface GeneratePageMetadataOptions {
  title?: string;
  description?: string;
  keywords?: string;
  ogImage?: string;
  canonical?: string;
  noIndex?: boolean;
}

/**
 * 为页面生成 SEO metadata
 * @param config 网站配置
 * @param options 页面特定的 SEO 选项
 * @returns Metadata 对象
 */
export function generatePageMetadata(
  config: IndexConfigVO,
  options: GeneratePageMetadataOptions = {}
): Metadata {
  const seoConfig = config.seo_meta_config;
  const websiteConfig = config.website_config;
  
  const title = options.title 
    ? `${options.title} - ${websiteConfig.website_name}`
    : seoConfig.title || websiteConfig.website_name;
    
  const description = options.description || seoConfig.description;
  const keywords = options.keywords || seoConfig.keywords;
  const ogImage = options.ogImage || seoConfig.og_image;
  
  return {
    title,
    description,
    keywords,
    authors: [{ name: seoConfig.author || websiteConfig.website_owner }],
    robots: options.noIndex ? 'noindex, nofollow' : (seoConfig.robots || 'index, follow'),
    openGraph: {
      title: options.title || seoConfig.og_title || title,
      description,
      images: ogImage ? [{ url: ogImage }] : undefined,
      siteName: websiteConfig.website_name,
      type: 'website',
    },
    twitter: {
      card: 'summary_large_image',
      title: options.title || seoConfig.og_title || title,
      description,
      images: ogImage ? [ogImage] : undefined,
    },
    alternates: {
      canonical: options.canonical,
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

/**
 * 生成文章页面的结构化数据
 */
export function generateArticleJsonLd(article: {
  title: string;
  description: string;
  publishedTime: string;
  modifiedTime?: string;
  author: string;
  authorImage?: string;
  image?: string;
  url: string;
  category?: string[];
  tags?: string[];
}) {
  return {
    "@context": "https://schema.org",
    "@type": "Article",
    "headline": article.title,
    "description": article.description,
    "image": article.image,
    "datePublished": article.publishedTime,
    "dateModified": article.modifiedTime || article.publishedTime,
    "author": {
      "@type": "Person",
      "name": article.author,
      "image": article.authorImage
    },
    "publisher": {
      "@type": "Organization",
      "name": article.author,
    },
    "mainEntityOfPage": {
      "@type": "WebPage",
      "@id": article.url
    },
    "keywords": article.tags?.join(', '),
    "articleSection": article.category?.join(', ')
  };
}
