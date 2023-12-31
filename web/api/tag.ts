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
export const getTagList = (apiDomain: string) => httpRequest.get(`${apiDomain}/${prefix}`)
export const getTagByRoute = (apiDomain: string, route: string) => httpRequest.get(`${apiDomain}/${prefix}/route/${route}`)
