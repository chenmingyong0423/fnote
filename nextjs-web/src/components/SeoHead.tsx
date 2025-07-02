'use client';

import { useEffect } from 'react';
import type { IndexConfigVO } from '../api/config';

interface SeoHeadProps {
  config: IndexConfigVO;
}

export default function SeoHead({ config }: SeoHeadProps) {
  useEffect(() => {
    // 动态设置网站图标
    if (config.website_config.website_icon) {
      const link = document.querySelector("link[rel*='icon']") || document.createElement('link');
      link.setAttribute('rel', 'shortcut icon');
      link.setAttribute('href', (process.env.NEXT_PUBLIC_SERVER_HOST || '') + config.website_config.website_icon);
      document.getElementsByTagName('head')[0].appendChild(link);
    }

    if (typeof window !== "undefined") {
      // 动态设置 canonical URL
      const canonical = document.querySelector("link[rel='canonical']") || document.createElement('link');
      canonical.setAttribute('rel', 'canonical');
      canonical.setAttribute('href', window.location.href);
      document.getElementsByTagName('head')[0].appendChild(canonical);

      // 添加 JSON-LD 结构化数据
      const jsonLd = {
        "@context": "https://schema.org",
        "@type": "WebSite",
        "name": config.website_config.website_name,
        "url": window.location.origin,
        "description": config.seo_meta_config.description,
        "author": {
          "@type": "Person",
          "name": config.website_config.website_owner,
          "image": config.website_config.website_owner_avatar,
          "description": config.website_config.website_owner_profile
        },
        "publisher": {
          "@type": "Organization",
          "name": config.website_config.website_name,
          "logo": {
            "@type": "ImageObject",
            "url": (process.env.NEXT_PUBLIC_SERVER_HOST || '') + config.website_config.website_icon
          }
        }
      };

      // 移除已存在的 JSON-LD 脚本
      const existingScript = document.querySelector('script[type="application/ld+json"]');
      if (existingScript) {
        existingScript.remove();
      }

      // 添加新的 JSON-LD 脚本
      const script = document.createElement('script');
      script.type = 'application/ld+json';
      script.textContent = JSON.stringify(jsonLd);
      document.head.appendChild(script);

      // 添加其他自定义 meta 标签
      const customMetas = [
        { name: 'theme-color', content: '#ffffff' },
        { name: 'msapplication-TileColor', content: '#da532c' },
        { property: 'og:type', content: 'website' },
        { property: 'og:locale', content: 'zh_CN' },
        { name: 'twitter:card', content: 'summary_large_image' },
      ];

      customMetas.forEach(({ name, property, content }) => {
        const selector = name ? `meta[name="${name}"]` : `meta[property="${property}"]`;
        let meta = document.querySelector(selector) as HTMLMetaElement;
        
        if (!meta) {
          meta = document.createElement('meta');
          if (name) meta.setAttribute('name', name);
          if (property) meta.setAttribute('property', property);
          document.head.appendChild(meta);
        }
        
        meta.setAttribute('content', content);
      });
    }
  }, [config]);

  return null;
}
