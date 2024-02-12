import type { IMenu } from "@/api/category";
import { defineStore } from "pinia";

export const useMenuStore = defineStore("menu", {
  state: () => ({
    menuList: [] as IMenu[], //菜单列表
  }),
  actions: {
    setMenuList(menuList: IMenu[]) {
      this.menuList = menuList;
    },
    getMenuList() {
      return this.menuList;
    },
  },
});
