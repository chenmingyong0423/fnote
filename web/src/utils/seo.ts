import type { Metadata } from "next";
import type { CommonConfigVO } from "@/src/api/config";
import type { PostDetail } from "@/src/api/posts";
import { resolvePublicUrl } from "@/src/utils/publicUrl";

const DEFAULT_SITE_URL = "http://localhost:3000";

const SIMPLE_ROBOTS_RULES = new Set([
  "all",
  "none",
  "index",
  "noindex",
  "follow",
  "nofollow",
  "nosnippet",
  "noimageindex",
  "notranslate",
  "nositelinkssearchbox",
]);

const VALUE_ROBOTS_RULES = [
  /^max-snippet:\s*-?\d+$/i,
  /^max-video-preview:\s*-?\d+$/i,
  /^max-image-preview:\s*(none|standard|large)$/i,
];

export function getSiteUrl() {
  const configuredUrl = process.env.BASE_HOST || DEFAULT_SITE_URL;

  try {
    return new URL(configuredUrl.endsWith("/") ? configuredUrl : `${configuredUrl}/`);
  } catch {
    return new URL(DEFAULT_SITE_URL);
  }
}

/**
 * 后台的 robots 字段历史上是自由文本，错误配置会原样输出到 meta。
 * 只保留搜索引擎支持的规则；完全无效时回到可索引的安全默认值。
 */
export function normalizeRobots(value?: string): NonNullable<Metadata["robots"]> {
  const rules = (value || "")
    .split(",")
    .map((rule) => rule.trim().toLowerCase())
    .filter(
      (rule) =>
        SIMPLE_ROBOTS_RULES.has(rule) ||
        VALUE_ROBOTS_RULES.some((pattern) => pattern.test(rule)),
    );

  return rules.length > 0 ? Array.from(new Set(rules)).join(", ") : "index, follow";
}

function uniqueNonEmpty(values: Array<string | undefined>) {
  return Array.from(
    new Set(values.map((value) => value?.trim()).filter((value): value is string => Boolean(value))),
  );
}

export function buildWebsiteJsonLd(config: CommonConfigVO) {
  const siteUrl = getSiteUrl();
  const siteName = config.website_meta.website_name || config.seo_meta.title;
  const owner = config.seo_meta.author || config.website_meta.website_owner;
  const keywordAliases = (config.seo_meta.keywords || "")
    .split(/[,，]/)
    .map((keyword) => keyword.trim())
    .filter(
      (keyword) =>
        keyword.length >= 2 &&
        keyword.length <= 40 &&
        Boolean(owner) &&
        keyword.includes(owner),
    );
  const alternateNames = uniqueNonEmpty([
    config.seo_meta.title,
    config.seo_meta.og_title,
    owner,
    ...keywordAliases,
    siteUrl.hostname,
  ]).filter((name) => name !== siteName);
  const ownerImage = resolvePublicUrl(config.website_meta.website_owner_avatar);

  return {
    "@context": "https://schema.org",
    "@type": "WebSite",
    "@id": `${siteUrl.toString()}#website`,
    url: siteUrl.toString(),
    name: siteName,
    ...(alternateNames.length > 0 ? { alternateName: alternateNames } : {}),
    description: config.seo_meta.description,
    inLanguage: "zh-CN",
    publisher: {
      "@type": "Person",
      "@id": `${siteUrl.toString()}#person`,
      name: owner,
      description: config.website_meta.website_owner_profile,
      ...(ownerImage ? { image: ownerImage } : {}),
      url: siteUrl.toString(),
    },
  };
}

export function buildBlogPostingJsonLd(post: PostDetail, config: CommonConfigVO) {
  const siteUrl = getSiteUrl();
  const articleUrl = new URL(`/posts/${post._id}`, siteUrl).toString();
  const authorUrl = new URL("/about-me", siteUrl).toString();
  const siteName = config.website_meta.website_name || config.seo_meta.title;
  const owner = config.seo_meta.author || config.website_meta.website_owner;
  const articleImage = resolvePublicUrl(post.cover_img);
  const ownerImage = resolvePublicUrl(config.website_meta.website_owner_avatar);

  return {
    "@context": "https://schema.org",
    "@type": "BlogPosting",
    "@id": `${articleUrl}#article`,
    mainEntityOfPage: {
      "@type": "WebPage",
      "@id": articleUrl,
    },
    headline: post.title,
    description: post.meta_description || post.summary,
    ...(articleImage ? { image: [articleImage] } : {}),
    datePublished: new Date(post.created_at * 1000).toISOString(),
    dateModified: new Date((post.updated_at || post.created_at) * 1000).toISOString(),
    author: {
      "@type": "Person",
      name: post.author,
      url: authorUrl,
    },
    publisher: {
      "@type": "Person",
      "@id": `${siteUrl.toString()}#person`,
      name: owner,
      url: authorUrl,
      ...(ownerImage ? { image: ownerImage } : {}),
    },
    isPartOf: {
      "@type": "WebSite",
      "@id": `${siteUrl.toString()}#website`,
      name: siteName,
      url: siteUrl.toString(),
    },
    inLanguage: "zh-CN",
    ...(post.category.length > 0
      ? { articleSection: post.category.map((category) => category.name) }
      : {}),
    ...(post.tags.length > 0 ? { keywords: post.tags.map((tag) => tag.name) } : {}),
  };
}

export function serializeJsonLd(value: unknown) {
  return JSON.stringify(value).replace(/</g, "\\u003c");
}
