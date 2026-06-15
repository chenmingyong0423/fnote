import { create } from "zustand";

interface ThemeState {
  isDark: boolean;
  toggleDark: () => void;
  setDark: (val: boolean) => void;
}

function getInitialDark() {
  if (typeof window === "undefined") return false;

  const saved = localStorage.getItem("theme-dark");
  if (saved === "1") return true;
  if (saved === "0") return false;

  return window.matchMedia?.("(prefers-color-scheme: dark)").matches ?? false;
}

export const useThemeStore = create<ThemeState>((set, get) => ({
  isDark: getInitialDark(),
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
