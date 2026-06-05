import { create } from "zustand";

interface ThemeState {
  isDark: boolean;
  toggleDark: () => void;
  setDark: (val: boolean) => void;
}

export const useThemeStore = create<ThemeState>((set, get) => ({
  isDark: false,
  toggleDark: () => {
    const next = !get().isDark;
    set({ isDark: next });
    if (typeof window !== "undefined") {
      document.documentElement.classList.toggle("dark", next);
      localStorage.setItem("theme-dark", next ? "1" : "0");
    }
  },
  setDark: (val: boolean) => {
    set({ isDark: val });
    if (typeof window !== "undefined") {
      document.documentElement.classList.toggle("dark", val);
      localStorage.setItem("theme-dark", val ? "1" : "0");
    }
  },
}));

// 初始化时自动读取 localStorage
if (typeof window !== "undefined") {
  const saved = localStorage.getItem("theme-dark");
  if (saved === "1") {
    document.documentElement.classList.add("dark");
    useThemeStore.getState().setDark(true);
  } else if (saved === "0") {
    document.documentElement.classList.remove("dark");
    useThemeStore.getState().setDark(false);
  }
}
