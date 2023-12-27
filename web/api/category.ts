import httpRequest from "./http";

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

const prefix = ""

export const getMenus = () => httpRequest.get(prefix + "/menus")
export const getCategoriesAndTags = () => httpRequest.get(prefix + "/categories")

export const getCategoryByRoute = (route: string) => httpRequest.get(prefix + `/categories/route/${route}`)


