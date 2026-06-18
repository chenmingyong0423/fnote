import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark'

const THEME_STORAGE_KEY = 'fnote-admin-theme'

const getInitialTheme = (): ThemeMode => {
  const savedTheme = localStorage.getItem(THEME_STORAGE_KEY)
  if (savedTheme === 'light' || savedTheme === 'dark') {
    return savedTheme
  }

  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>(getInitialTheme())

  const applyTheme = (themeMode: ThemeMode) => {
    const isDark = themeMode === 'dark'
    document.documentElement.classList.toggle('theme-dark', isDark)
    document.documentElement.dataset.theme = themeMode
    document.documentElement.style.colorScheme = themeMode
    localStorage.setItem(THEME_STORAGE_KEY, themeMode)
  }

  const toggleTheme = () => {
    mode.value = mode.value === 'dark' ? 'light' : 'dark'
  }

  watch(mode, applyTheme, { immediate: true })

  return { mode, toggleTheme }
})
