'use client';

import { useEffect, useState, useCallback } from 'react';
import { useConfigStore, useWebsiteInfo } from '../store/configStore';
import { collectVisitLog, VisitLogRequest } from '../api/statiscs';

// 用于跟踪API调用时间的变量
const API_CACHE = {
  visitLog: 0
};

// 设定API调用最小间隔时间（毫秒）
const MIN_CALL_INTERVAL = 60000; // 1分钟

export function useLayoutSetup() {
  // 使用 Zustand store 替代 Context
  const config = useConfigStore();
  // 使用选择器获取网站信息，避免不必要的重渲染
  const websiteInfo = useWebsiteInfo();
  
  // 收集访问日志
  const collectVisit = useCallback(async () => {
    try {
      // 检查是否需要调用API
      const now = Date.now();
      if (now - API_CACHE.visitLog < MIN_CALL_INTERVAL) {
        console.log('Skip collectVisit, called recently');
        return;
      }
      
      // 更新API调用时间
      API_CACHE.visitLog = now;
      
      const req: { url: string } = {
        url: window.location.href,
      };
      await collectVisitLog(req);
    } catch (error) {
      console.error('收集访问日志失败:', error);
    }
  }, []);
  
  // 滚动到顶部相关逻辑
  const [showScrollTop, setShowScrollTop] = useState(false);
  
  const handleScroll = useCallback(() => {
    if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 20) {
      setShowScrollTop(true);
    } else {
      setShowScrollTop(false);
    }
  }, []);
  
  const scrollToTop = () => {
    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    });
  };
  
  // 获取 JSON-LD 结构化数据
  const getJsonLd = useCallback(() => {
    // 增加空值检查，避免访问null对象的属性
    const siteName = websiteInfo?.website_config?.website_name || 'fnote';
    const siteURL = process.env.NEXT_PUBLIC_DOMAIN || '';
    
    return JSON.stringify({
      "@context": "https://schema.org",
      "@type": "WebSite",
      "name": siteName,
      "url": siteURL,
    });
  }, [websiteInfo]);
  
  // 初始化函数
  const initLayout = useCallback(async () => {
    console.log('Initializing layout...');
    
    // 初始化暗黑模式 (现在由 useInitConfig 处理)
    config.initDarkMode();

    // 获取网站统计数据 (现在由 useInitConfig 处理)
    await config.fetchWebsiteStats();

    // 获取网站信息 (现在由 useInitConfig 处理)
    await config.fetchWebsiteInfo();
    
    // 收集访问日志
    await collectVisit();
    
    // 添加滚动事件监听
    window.addEventListener('scroll', handleScroll);
    
    return () => {
      // 清理事件监听
      window.removeEventListener('scroll', handleScroll);
    };
  }, [config, collectVisit, handleScroll]);
  
  return {
    initLayout,
    collectVisit,
    scrollToTop,
    showScrollTop,
    getJsonLd,
  };
}