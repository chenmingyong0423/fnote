'use client';

import { useEffect } from 'react';
import { useLayoutSetup } from '../hooks/useLayoutSetup';
import useInitConfig from '../hooks/useInitConfig';
import { useConfigStore } from '../store/configStore';

// 用于跟踪是否已经初始化的全局变量
let isAppInitialized = false;

// 创建初始化组件，不渲染任何内容，只执行副作用
const AppInitializer = () => {
  // 使用Zustand store的初始化钩子
  const { errors } = useInitConfig();
  
  // 直接从Zustand store获取状态和方法
  const { isLoading } = useConfigStore();
  
  // 获取布局设置钩子
  const { initLayout } = useLayoutSetup();

  // 初始化应用程序
  useEffect(() => {
    // 如果已经初始化过，则不再重复初始化
    if (isAppInitialized) {
      return;
    }

    const init = async () => {
      try {
        // 标记为已初始化
        isAppInitialized = true;
        
        // 初始化布局相关功能
        const cleanup = await initLayout();
        
        // 返回清理函数
        return () => {
          if (typeof cleanup === 'function') {
            cleanup();
            // 组件卸载时重置初始化标志
            isAppInitialized = false;
          }
        };
      } catch (error) {
        console.error('Application initialization failed:', error);
        isAppInitialized = false;
      }
    };
    
    init();
  }, [initLayout]);
  
  // 监听错误和加载状态，可以在这里添加全局错误处理
  useEffect(() => {
    if (errors && (errors.websiteInfo || errors.websiteStats)) {
      console.warn('Config loading errors:', errors);
      // 这里可以添加错误通知或重试逻辑
    }
  }, [errors]);
  
  // 组件不渲染任何内容
  return null;
};

export default AppInitializer;