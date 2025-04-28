'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { WebsiteCountStats } from '../api/statiscs';
import { useConfigStore, useDarkMode } from '../store/configStore';

// 保留接口定义用于组件使用，但数据将从Zustand获取
interface WebsiteConfig {
  website_name: string;
  website_icon: string;
  website_owner?: string;
  website_owner_profile?: string;
  website_owner_avatar?: string;
  website_owner_email?: string;
  website_runtime?: number;
  website_records?: string[];
  [key: string]: any;
}

interface NoticeConfig {
  show: boolean;
  content: string;
  title?: string;
  publish_time?: number;
  [key: string]: any;
}

interface SocialInfoItem {
  name: string;
  link: string;
  icon: string;
}

interface PayInfoItem {
  name: string;
  qrcode: string;
}

interface SeoMetaConfig {
  title: string;
  description: string;
  keywords: string;
  author: string;
  robots: string;
  og_title: string;
  og_image: string;
}

interface ThirdPartyVerification {
  key: string;
  value: string;
  description?: string;
}

interface ConfigContextType {
  website_info: WebsiteConfig;
  notice_info: NoticeConfig;
  social_info_list: SocialInfoItem[];
  pay_info: PayInfoItem[];
  seo_meta_config: SeoMetaConfig;
  tpsv_list: ThirdPartyVerification[];
  website_count_stats: WebsiteCountStats;
  isDarkMode: boolean;
  toggleDarkMode: () => void;
  updateWebsiteStats: (stats: WebsiteCountStats) => void;
  metaVerificationList: { name: string; content: string }[];
}

// 默认值现在只作为类型兼容的临时变量，实际数据将从Zustand获取
const defaultConfig: ConfigContextType = {
  website_info: {
    website_name: 'fnote',
    website_icon: ''
  },
  notice_info: {
    show: true,
    content: '',
    title: ''
  },
  social_info_list: [],
  pay_info: [],
  seo_meta_config: {
    title: '',
    description: '',
    keywords: '',
    author: '',
    robots: '',
    og_title: '',
    og_image: ''
  },
  tpsv_list: [],
  website_count_stats: {
    post_count: 0,
    category_count: 0,
    tag_count: 0,
    comment_count: 0,
    like_count: 0,
    website_view_count: 0
  },
  isDarkMode: false,
  toggleDarkMode: () => {},
  updateWebsiteStats: () => {},
  metaVerificationList: []
};

// 仍然保留上下文对象，用于向后兼容，但内部使用Zustand
const ConfigContext = createContext<ConfigContextType>(defaultConfig);

// 导出当前上下文的使用钩子，在完全迁移到Zustand前保持兼容
export const useConfigContext = () => {
  const contextValue = useContext(ConfigContext);
  
  // 由于我们在过渡期，建议使用新的Zustand钩子
  console.warn('useConfigContext is deprecated, please use Zustand hooks from configStore.ts instead');
  
  return contextValue;
};

// 提供者现在只是一个包装，将Zustand状态注入到React上下文中
export const ConfigProvider: React.FC<{children: ReactNode}> = ({ children }) => {
  // 从Zustand获取状态和操作
  const store = useConfigStore();
  const isDarkMode = useDarkMode();
  
  // 创建与原上下文兼容的值对象
  const value: ConfigContextType = {
    // 从Zustand获取的数据与原上下文格式兼容
    website_info: store.website_info?.website_config || defaultConfig.website_info,
    notice_info: store.website_info?.notice_config ? {
      ...store.website_info.notice_config,
      show: !!store.website_info.notice_config.content
    } : defaultConfig.notice_info,
    social_info_list: store.website_info?.social_info_config?.social_info_list?.map(item => ({
      name: item.social_name,
      link: item.social_value,
      icon: item.css_class
    })) || [],
    pay_info: store.website_info?.pay_info_config?.map(item => ({
      name: item.name,
      qrcode: item.image
    })) || [],
    seo_meta_config: store.website_info?.seo_meta_config || defaultConfig.seo_meta_config,
    tpsv_list: store.website_info?.third_party_site_verification || [],
    website_count_stats: store.website_stats || defaultConfig.website_count_stats,
    isDarkMode,
    
    // 使用Zustand操作
    toggleDarkMode: store.toggleDarkMode,
    
    // 更新网站统计数据的函数，转发到Zustand
    updateWebsiteStats: (stats: WebsiteCountStats) => {
      store.setWebsiteStats(stats);
    },
    
    // 第三方站点验证元数据
    metaVerificationList: store.website_info?.third_party_site_verification?.map(item => ({
      name: item.key,
      content: item.value
    })) || []
  };
  
  return (
    <ConfigContext.Provider value={value}>
      {children}
    </ConfigContext.Provider>
  );
};