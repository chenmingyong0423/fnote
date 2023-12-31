import httpRequest from "./http";

// type TagsWithCountVO struct {
// 	Name  string `json:"name"`
// 	Count int64  `json:"count"`
// }
export interface ITagWithCount {
    name: string;
    route: string;
    count: number;
}

export interface ITagName {
    name: string;
}

const prefix = "tags"
export const getTagList = (apiBaseUrl: string) => httpRequest.get(`${apiBaseUrl}/${prefix}`)
export const getTagByRoute = (apiBaseUrl: string, route: string) => httpRequest.get(`${apiBaseUrl}/${prefix}/route/${route}`)
