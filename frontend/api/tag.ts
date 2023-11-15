import httpRequest from "./http";
export const getTagList = (name: string) => httpRequest.get(`/categories/categories/${name}/tags`)