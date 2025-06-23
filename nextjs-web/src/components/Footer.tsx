import React from "react";

const Footer: React.FC = () => (
  <footer className="w-full bg-white border-t border-gray-200 mt-8 h-[60px] px-4 mx-auto max-w-7xl grid grid-cols-12 items-center">
    {/* 左留白 */}
    <div className="col-span-1 lg:col-span-2" />
    {/* 主体内容 10/12 或 8/12 居中 */}
    <div className="col-span-10 lg:col-span-8 flex items-center justify-center w-full h-full text-center text-gray-500 text-sm">
      © {new Date().getFullYear()} Copyright © 2024 - Designed by{" "}
      <a
        href="https://github.com/chenmingyong0423/fnote"
        className="ml-1 text-blue-600 hover:underline"
      >
        Fnote
      </a>
    </div>
    {/* 右留白 */}
    <div className="col-span-1 lg:col-span-2" />
  </footer>
);

export default Footer;
