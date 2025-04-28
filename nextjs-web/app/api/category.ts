import httpRequest from "../utils/http";

export interface IMenu {
  name: string;
  route: string;
}

export interface ICategoryWithCount {
  name: string;
  route: string;
  description: string;
  count: number;
}

export interface ICategoryName {
  name: string;
}

const prefix = "/categories";

export const getMenus = () => httpRequest.get(`/menus`);
export const getCategoriesAndTags = () => httpRequest.get(prefix);
export const getCategoryByRoute = (route: string) =>
  httpRequest.get(`${prefix}/route/${route}`);