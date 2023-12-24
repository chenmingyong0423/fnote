import httpRequest from "./http";

export interface IMenu {
    name: string;
    route: string;
}

export interface ICategoryAndTags {
    categories: ICategory[];
    tags: ITag[];
}

export interface ICategory {
    name: string;
    route: string;
    description: string;
    count: number;
}

export interface ITag {
    name: string;
    count: number;
}

const prefix = ""

export const getMenus = () => httpRequest.get(prefix + "/menus")
export const getCategoriesAndTags = () => httpRequest.get(prefix + "/categories")


