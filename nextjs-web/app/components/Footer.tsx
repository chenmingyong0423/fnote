'use client';

import { useConfigContext } from "../context/ConfigContext";

const Footer = () => {
  const { website_info } = useConfigContext();
  
  // 计算年份，基于网站运行时间
  const date = website_info.website_runtime 
    ? new Date(website_info.website_runtime * 1000)
    : new Date();
  const currentYear = date.getFullYear();
  
  // 备案信息，从配置中获取
  const records = website_info.website_records || [];

  return (
    <footer className="flex flex-col justify-evenly items-center p-10 dark:text-gray-300 dark:bg-gray-900 shadow-xl">
      <div>
        Copyright © {currentYear} - Designed by
        <a
          className="no-underline text-black font-bold hover:underline dark:text-white ml-1"
          href="https://github.com/chenmingyong0423/fnote"
          target="_blank"
          rel="noopener noreferrer"
        >
          Fnote
        </a>
      </div>
      <div className="flex items-center lg:flex-row flex-col">
        {records.map((item, index) => (
          <span key={index} dangerouslySetInnerHTML={{ __html: item }} />
        ))}
      </div>
      <div></div>
    </footer>
  );
};

export default Footer;