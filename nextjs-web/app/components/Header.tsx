'use client';

import Link from 'next/link';
import Image from 'next/image';
import { useRef, useEffect, useState } from 'react';
import { useConfigContext } from '../context/ConfigContext';
import Menu from './Menu';
import { IMenu, getMenus } from '../api/category';
import { IResponse, IListData } from '../utils/http';

interface HeaderProps {
  className?: string;
}

const Header = ({ className = '' }: HeaderProps) => {
  const { isBlackMode, toggleDarkMode, website_info } = useConfigContext();
  const [menuList, setMenuList] = useState<IMenu[]>([]);
  const [scrollCount, setScrollCount] = useState(0);
  const [showMobileMenu, setShowMobileMenu] = useState(false);
  const headerRef = useRef<HTMLDivElement>(null);
  
  // 获取菜单数据
  useEffect(() => {
    // 默认菜单，如果API请求失败时使用
    const defaultMenus: IMenu[] = [
      { name: '后端', route: 'backend' },
      { name: '前端', route: 'frontend' },
      { name: '人工智能', route: 'ai' },
      { name: '工具', route: 'tools' },
    ];
    
    // 使用 category.ts 中定义的 getMenus API 函数获取菜单数据
    const fetchMenuList = async () => {
      try {
        const response = await getMenus();
        if (response && response.data) {
          const menuData = response.data as IListData<IMenu>;
          setMenuList(menuData.list || []);
        } else {
          setMenuList(defaultMenus);
        }
      } catch (error) {
        console.error('获取菜单数据失败:', error);
        setMenuList(defaultMenus);
      }
    };

    fetchMenuList();
  }, []);

  // 处理滚动隐藏/显示导航栏的效果
  useEffect(() => {
    const headerScroll = () => {
      if (!headerRef.current) return;
      
      if (document.documentElement.scrollTop > scrollCount) {
        headerRef.current.style.top = `-${headerRef.current.clientHeight}px`;
      } else {
        headerRef.current.style.top = '0px';
      }
      
      setScrollCount(document.documentElement.scrollTop);
    };
    
    window.addEventListener('scroll', headerScroll);
    return () => window.removeEventListener('scroll', headerScroll);
  }, [scrollCount]);

  return (
    <>
      <header 
        ref={headerRef}
        className={`bg-white backdrop-blur-xl fixed top-0 w-full z-50 flex justify-between items-center p-1 dark:bg-gray-900 duration-200 ease-linear max-h-[70px] select-none ${className}`}
      >
        <div>
          <Link href="/">
            <Image 
              src={website_info.website_icon || "/vercel.svg"}
              alt="Logo"
              width={60}
              height={60}
              className="w-15 h-15 rounded-full mx-5 cursor-pointer hover:rotate-[360deg] ease-out duration-1000 sm:mr-0 select-none"
            />
          </Link>
        </div>
        
        <div className="bg-transparent md:block hidden">
          <Menu items={menuList} />
        </div>
        
        <div className="ml-auto pr-5 flex gap-x-4">
          {/* 暗黑模式切换按钮 */}
          <button 
            onClick={toggleDarkMode} 
            className="cursor-pointer text-2.5rem text-[#86909c] dark:text-gray-300 dark:hover:text-white"
            aria-label={isBlackMode ? "切换到亮色模式" : "切换到暗色模式"}
          >
            {isBlackMode ? (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            ) : (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
            )}
          </button>

          {/* 搜索按钮 */}
          <Link 
            href="/search?keyword=" 
            className="cursor-pointer text-2.5rem text-[#86909c] dark:text-gray-300 dark:hover:text-white"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </Link>

          {/* 小屏幕菜单按钮 */}
          <button 
            className="cursor-pointer text-2.5rem text-[#86909c] dark:text-gray-300 dark:hover:text-white md:hidden"
            onClick={() => setShowMobileMenu(!showMobileMenu)}
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
      </header>
      
      {/* 移动端菜单 */}
      {showMobileMenu && (
        <div className="fixed top-0 right-0 w-full h-full z-[100] md:hidden">
          {/* 移动端菜单背景遮罩 */}
          <div 
            className="absolute inset-0 bg-black/50"
            onClick={() => setShowMobileMenu(false)}
          ></div>
          
          {/* 移动端菜单内容 */}
          <div className="absolute top-0 right-0 w-4/5 h-full bg-white dark:bg-gray-900 p-4 overflow-y-auto">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-bold dark:text-white">菜单</h2>
              <button 
                onClick={() => setShowMobileMenu(false)}
                className="p-2 text-gray-600 dark:text-gray-300"
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <Menu items={menuList} />
          </div>
        </div>
      )}
    </>
  );
};

export default Header;