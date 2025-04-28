'use client';

import { useRef } from 'react';
import Header from './Header';
import Footer from './Footer';
import MyToast from './MyToast';
import ScrollToTop from './ScrollToTop';
import { useConfigContext } from '../context/ConfigContext';
import { useLayoutSetup } from '../hooks/useLayoutSetup';
import SEO from './SEO';
import AppInitializer from './AppInitializer';

export default function BlogLayout({ children }: { children: React.ReactNode }) {
  const headerRef = useRef<HTMLDivElement>(null);
  const config = useConfigContext();
  
  // 使用我们创建的 layout setup hook
  const {
    showScrollTop,
    scrollToTop
  } = useLayoutSetup();

  return (
    <div className="dark:bg-[#03080c]">
      {/* SEO 和元数据组件 */}
      <SEO />
      
      {/* 全局初始化组件（收集访问日志、初始化全局功能等） */}
      <AppInitializer />
      
      <MyToast />
      
      <div ref={headerRef}>
        <Header className="slide-down" />
      </div>
      
      <div className="bg-[#F0F2F5] dark:bg-[#03080c] pt-25 p-5">
        <div className="w-[90%] m-auto sm:w-[99%]">
          {children}
        </div>
      </div>
      
      <div>
        <Footer />
      </div>
      
      {/* 使用 showScrollTop 状态控制回到顶部按钮的显示 */}
      <ScrollToTop show={showScrollTop} onClick={scrollToTop} />
      
      <style jsx global>{`
        @keyframes slideDown {
          0% {
            transform: translateY(-100%);
          }
          100% {
            transform: translateY(0);
          }
        }

        .slide-down {
          animation: slideDown 1s ease;
          animation-iteration-count: 1;
        }
      `}</style>
    </div>
  );
}