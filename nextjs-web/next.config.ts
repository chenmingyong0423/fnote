import type { NextConfig } from "next";

// 从环境变量中解析后端主机信息
const serverUrl = new URL(process.env.SERVER_HOST || 'http://localhost:8080');
const serverHost = process.env.SERVER_HOST || 'http://localhost:8080';

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${serverHost}/:path*`, // 转发到后端服务器
      },
      {
        source: '/static/:path*',
        destination: `${serverHost}/static/:path*`, // 转发到后端服务器的静态资源
      },
    ];
  },
  
  // 配置允许的图片域名和路径前缀
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