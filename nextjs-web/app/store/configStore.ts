import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { getWebsiteInfo, IWebsiteInfo } from '../api/config';
import { getWebsiteCountStats, WebsiteCountStats } from '../api/statiscs';

// 配置状态接口
interface ConfigState {
  // 网站信息
  website_info: IWebsiteInfo | null;
  // 网站统计数据
  website_stats: WebsiteCountStats | null;
  // 主题设置
  isDarkMode: boolean;
  // 是否正在加载数据
  isLoading: { [key: string]: boolean };
  // 错误信息
  error: { [key: string]: string | null };
  // 设置网站信息
  setWebsiteInfo: (info: IWebsiteInfo) => void;
  // 设置网站统计数据
  setWebsiteStats: (stats: WebsiteCountStats) => void;
  // 切换暗黑模式
  toggleDarkMode: () => void;
  // 初始化暗黑模式（从系统或存储的偏好中）
  initDarkMode: () => void;
  // 获取网站信息
  fetchWebsiteInfo: () => Promise<IWebsiteInfo | null>;
  // 获取网站统计数据
  fetchWebsiteStats: () => Promise<WebsiteCountStats | null>;
  // 设置加载状态
  setLoading: (key: string, isLoading: boolean) => void;
  // 设置错误信息
  setError: (key: string, error: string | null) => void;
}

// 创建配置store
export const useConfigStore = create<ConfigState>()(
  devtools(
    persist(
      (set, get) => ({
        // 初始状态
        website_info: null,
        website_stats: null,
        isDarkMode: false,
        isLoading: { websiteInfo: false, websiteStats: false },
        error: { websiteInfo: null, websiteStats: null },

        // 设置网站信息
        setWebsiteInfo: (info) => set({ website_info: info }),

        // 设置网站统计数据
        setWebsiteStats: (stats) => set({ website_stats: stats }),

        // 设置加载状态
        setLoading: (key, isLoading) => set((state) => ({
          isLoading: { ...state.isLoading, [key]: isLoading }
        })),

        // 设置错误信息
        setError: (key, error) => set((state) => ({
          error: { ...state.error, [key]: error }
        })),

        // 切换暗黑模式
        toggleDarkMode: () => {
          const newDarkMode = !get().isDarkMode;
          set({ isDarkMode: newDarkMode });
          
          // 更新文档根元素的类，用于CSS样式切换
          if (typeof document !== 'undefined') {
            if (newDarkMode) {
              document.documentElement.classList.add('dark');
            } else {
              document.documentElement.classList.remove('dark');
            }
          }
        },

        // 初始化暗黑模式
        initDarkMode: () => {
          // 检查用户偏好或系统设置
          const storedDarkMode = get().isDarkMode;
          const prefersDark = 
            typeof window !== 'undefined' && 
            window.matchMedia('(prefers-color-scheme: dark)').matches;
          
          // 如果已有存储的偏好，使用它；否则使用系统偏好
          const shouldBeDark = storedDarkMode !== undefined ? storedDarkMode : prefersDark;
          
          set({ isDarkMode: shouldBeDark });
          
          // 更新文档根元素的类，用于CSS样式切换
          if (typeof document !== 'undefined') {
            if (shouldBeDark) {
              document.documentElement.classList.add('dark');
            } else {
              document.documentElement.classList.remove('dark');
            }
          }
        },

        // 获取网站信息 - 使用已有的API接口函数
        fetchWebsiteInfo: async () => {
          const { setLoading, setError } = get();
          setLoading('websiteInfo', true);
          setError('websiteInfo', null);
          
          try {
            const response = await getWebsiteInfo();
            if (response && response.data) {
              set({ website_info: response.data });
              return response.data;
            }
            return null;
          } catch (error) {
            console.error('Failed to fetch website info:', error);
            setError('websiteInfo', error instanceof Error ? error.message : '获取网站信息失败');
            return null;
          } finally {
            setLoading('websiteInfo', false);
          }
        },

        // 获取网站统计数据 - 使用已有的API接口函数
        fetchWebsiteStats: async () => {
          const { setLoading, setError } = get();
          setLoading('websiteStats', true);
          setError('websiteStats', null);
          
          try {
            const response = await getWebsiteCountStats();
            if (response && response.data) {
              set({ website_stats: response.data });
              return response.data;
            }
            return null;
          } catch (error) {
            console.error('Failed to fetch website statistics:', error);
            setError('websiteStats', error instanceof Error ? error.message : '获取网站统计数据失败');
            return null;
          } finally {
            setLoading('websiteStats', false);
          }
        },
      }),
      {
        name: 'fnote-config-storage', // 存储的键名
        partialize: (state) => ({
          isDarkMode: state.isDarkMode, // 只将isDarkMode持久化到本地存储
        }),
      }
    )
  )
);

// 导出一些常用的选择器函数，便于组件使用
export const useWebsiteInfo = () => useConfigStore((state) => state.website_info);
export const useWebsiteStats = () => useConfigStore((state) => state.website_stats);
export const useDarkMode = () => useConfigStore((state) => state.isDarkMode);
export const useConfigLoading = () => useConfigStore((state) => state.isLoading);
export const useConfigError = () => useConfigStore((state) => state.error);