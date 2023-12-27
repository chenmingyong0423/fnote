import httpRequest from "./http";

// type TagsWithCountVO struct {
// 	Name  string `json:"name"`
// 	Count int64  `json:"count"`
// }
export interface ITagWithCount {
    name:string;
    route:string;
    count:number;
}
export const getTagList = () => httpRequest.get(`/tags`)