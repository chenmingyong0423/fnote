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

const prefix = "categories"

export const getMenus = (apiBaseUrl: string) => httpRequest.get(`${apiBaseUrl}/menus`)
export const getCategoriesAndTags = (apiBaseUrl: string) => httpRequest.get(`${apiBaseUrl}/${prefix}`)

export const getCategoryByRoute = (apiBaseUrl: string, route: string) => httpRequest.get(`${apiBaseUrl}/${prefix}/route/${route}`)


