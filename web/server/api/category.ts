import httpRequest from "./http";

export interface IMenu {
    name: string;
    route: string;
}

const prefix = ""

export const getMenus = () => httpRequest.get(prefix + "/menus")


