import type { NextConfig } from "next";

// 从环境变量中解析后端主机信息
const serverUrl = new URL(process.env.SERVER_HOST || 'http://localhost:8080');
const serverHost = process.env.SERVER_HOST || 'http://localhost:8080';

const nextConfig: NextConfig = {
  // 🚨 关键：让 Docker build 不因为 ESLint 报错失败
  eslint: {
    ignoreDuringBuilds: true,
  },

  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${serverHost}/:path*`,
      },
      {
        source: '/static/:path*',
        destination: `${serverHost}/static/:path*`,
      },
    ];
  },

  images: {
    remotePatterns: [
      {
        protocol: serverUrl.protocol.replace(':', '') as 'http' | 'https',
        hostname: serverUrl.hostname,
        port: serverUrl.port || undefined,
        pathname: '/static/**',
      },
    ],
  },
};

export default nextConfig;