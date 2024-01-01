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

const prefix = "/tags"
export const getTagList = () => httpRequest.get(`$${prefix}`)
export const getTagByRoute = (route: string) => httpRequest.get(`${prefix}/route/${route}`)
