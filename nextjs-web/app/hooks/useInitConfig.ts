import { useEffect } from 'react';
import { useConfigStore, useWebsiteStats, useConfigLoading, useConfigError } from '../store/configStore';

/**
 * 初始化全局配置的Hook
 * 该Hook会在应用初始化时获取网站配置信息和统计数据
 * 使用Zustand进行状态管理
 */
export default function useInitConfig() {
  // 从Zustand store获取状态和方法
  const { 
    fetchWebsiteInfo, 
    fetchWebsiteStats,
    initDarkMode
  } = useConfigStore();
  
  // 使用选择器获取特定状态，避免不必要的重渲染
  const website_stats = useWebsiteStats();
  const isLoading = useConfigLoading();
  const errors = useConfigError();

  // 初始化网站配置信息和暗黑模式
  useEffect(() => {
    // 初始化暗黑模式
    initDarkMode();
    
    // 获取网站配置信息
    fetchWebsiteInfo();
    
    // 获取网站统计数据
    fetchWebsiteStats();
    
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  
  // 定期更新网站统计数据
  useEffect(() => {
    // 定时刷新网站统计数据，每5分钟获取一次
    const statsInterval = setInterval(() => {
      fetchWebsiteStats();
    }, 5 * 60 * 1000);

    return () => {
      clearInterval(statsInterval);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { 
    website_stats, 
    isLoading,
    errors,
    refreshStats: fetchWebsiteStats,
    refreshInfo: fetchWebsiteInfo
  };
}