import { create } from 'zustand';
import type { IndexConfigVO } from '../api/config';

interface ConfigState {
  config: IndexConfigVO | null;
  setConfig: (config: IndexConfigVO) => void;
  clearConfig: () => void;
}

export const useConfigStore = create<ConfigState>((set) => ({
  config: null,
  setConfig: (config) => set({ config }),
  clearConfig: () => set({ config: null }),
}));
