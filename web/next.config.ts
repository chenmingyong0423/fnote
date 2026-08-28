import type { NextConfig } from "next";

// 从环境变量中解析后端主机信息
const serverUrl = new URL(process.env.SERVER_HOST || "http://localhost:8080");
const serverHost = process.env.SERVER_HOST || "http://localhost:8080";

const nextConfig: NextConfig = {
  async redirects() {
    return [
      {
        source: "/about",
        destination: "/about-me",
        permanent: true,
      },
    ];
  },

  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${serverHost}/:path*`,
      },
      {
        source: "/static/:path*",
        destination: `${serverHost}/static/:path*`,
      },
      {
        source: "/sitemap.xml",
        destination: `${serverHost}/static/sitemap.xml`,
      },
      {
        source: "/robots.txt",
        destination: `${serverHost}/static/robots.txt`,
      },
    ];
  },

  images: {
    remotePatterns: [
      {
        protocol: serverUrl.protocol.replace(":", "") as "http" | "https",
        hostname: serverUrl.hostname,
        port: serverUrl.port || undefined,
        pathname: "/static/**",
      },
    ],
  },
};

export default nextConfig;
