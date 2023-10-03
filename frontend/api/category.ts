import httpRequest from "./http";

export interface IMenu {
    name: string;
    route: string;
}

const prefix = "/categories"

const getMenus = () => {
    return httpRequest.get(prefix + "/menus")
};

export {
    getMenus
}