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

export const getMenus = (apiDomain: string) => httpRequest.get(`${apiDomain}/menus`)
export const getCategoriesAndTags = (apiDomain: string) => httpRequest.get(`${apiDomain}/${prefix}`)

export const getCategoryByRoute = (apiDomain: string, route: string) => httpRequest.get(`${apiDomain}/${prefix}/route/${route}`)


