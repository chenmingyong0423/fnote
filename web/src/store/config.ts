import { create } from 'zustand';
import type { CommonConfigVO } from '../api/config';

interface CommonConfigState {
  config: CommonConfigVO | null;
  setConfig: (config: CommonConfigVO) => void;
  clearConfig: () => void;
}

export const useConfigStore = create<CommonConfigState>((set) => ({
  config: null,
  setConfig: (config) => set({ config }),
  clearConfig: () => set({ config: null }),
}));
