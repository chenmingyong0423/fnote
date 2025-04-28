'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { IMenu } from '../api/category';

interface MenuProps {
  items: IMenu[];
}

const Menu: React.FC<MenuProps> = ({ items }) => {
  const pathname = usePathname();

  // 判断路径是否为当前活动菜单项
  const isActive = (path: string): boolean => {
    return pathname === path || 
           (path.startsWith('/categories/') && pathname.startsWith(path));
  };

  return (
    <nav>
      <div className="list-none flex gap-x-10 md:flex-row flex-col md:gap-y-0 gap-y-4 px-5">
        {/* 首页链接 */}
        <Link
          href="/"
          className={`flex items-center gap-x-4 cursor-pointer p-2 hover:border-b-2 hover:border-blue-500 ${
            pathname === '/' ? 'text-[#1e80ff]' : ''
          }`}
        >
          <span>
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
            </svg>
          </span>
          <span>主页</span>
        </Link>
        
        {/* 文章导航链接 */}
        <Link
          href="/navigation"
          className={`flex items-center gap-x-4 cursor-pointer p-2 hover:border-b-2 hover:border-blue-500 ${
            pathname === '/navigation' ? 'text-[#1e80ff]' : ''
          }`}
        >
          <span>
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clipRule="evenodd" />
            </svg>
          </span>
          <span>文章导航</span>
        </Link>
        
        {/* 动态分类菜单项 */}
        {items.map((item, index) => {
          const categoryPath = `/categories/${item.route}`;
          return (
            <Link
              key={index}
              href={categoryPath}
              className={`flex items-center gap-x-4 cursor-pointer p-2 hover:border-b-2 hover:border-blue-500 ${
                pathname === categoryPath || pathname.startsWith(categoryPath) ? 'text-[#1e80ff]' : ''
              }`}
            >
              <span>
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clipRule="evenodd" />
                </svg>
              </span>
              <span>{item.name}</span>
            </Link>
          );
        })}
        
        {/* 友链链接 */}
        <Link
          href="/friend"
          className={`flex items-center gap-x-4 cursor-pointer p-2 hover:border-b-2 hover:border-blue-500 ${
            pathname === '/friend' ? 'text-[#1e80ff]' : ''
          }`}
        >
          <span>
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
            </svg>
          </span>
          <span>友链</span>
        </Link>
        
        {/* 关于我链接 */}
        <Link
          href="/about-me"
          className={`flex items-center gap-x-4 cursor-pointer p-2 hover:border-b-2 hover:border-blue-500 ${
            pathname === '/about-me' ? 'text-[#1e80ff]' : ''
          }`}
        >
          <span>
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
            </svg>
          </span>
          <span>关于我</span>
        </Link>
      </div>
    </nav>
  );
};

export default Menu;