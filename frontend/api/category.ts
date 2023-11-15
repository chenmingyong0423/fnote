import httpRequest from "./http";

export interface IMenu {
    name: string;
    route: string;
}

const prefix = "/categories"

export const getMenus = () => httpRequest.get(prefix + "/menus")


