import { create } from 'zustand';

// 示例：全局 UI 状态
interface UIState {
  darkMode: boolean;
  toggleDarkMode: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  darkMode: false,
  toggleDarkMode: () => set((state) => ({ darkMode: !state.darkMode })),
}));

// 你可以在 src/store 目录下继续扩展其他模块的 store
