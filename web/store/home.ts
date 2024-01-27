import { defineStore } from "pinia";

export const useHomeStore = defineStore("home", {
  state: () => ({
    searchVisible: false, // 搜索弹窗状态
    isBlackMode: false,
    showSmallScreenMenu: false,
  }),
  // 持久化存储
  persist: process.client && {
    storage: localStorage, //存储模式：localStorage || sessionStorage
    paths: ["classification"], //要存储的数据
  },
});
