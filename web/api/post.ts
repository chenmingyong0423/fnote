import httpRequest from "./http";
import * as domain from "domain";

export interface IPost {
    sug: string;
    author: string;
    title: string;
    summary: string;
    cover_img: string;
    categories: string[];
    tags: string[];
    like_count: number;
    comment_count: number;
    visit_count: number;
    sticky_weight: number;
    create_time: number;
}

export type PageRequest = {
    pageNo: number;
    pageSize: number;
    sortField?: string;
    sortOrder?: string;
    keyword?: string;
    categories?: string[];
    tags?: string[];
}

export interface IPostDetail {
    sug: string;
    author: string;
    title: string;
    summary: string;
    cover_img: string;
    category: string[];
    tags: string[];
    like_count: number;
    comment_count: number;
    visit_count: number;
    sticky_weight: number;
    create_time: number;
    content: string;
    meta_description: string;
    meta_keywords: string;
    word_count: number;
    update_time: number;
    is_liked: boolean;
}

const prefix = "posts"
export const getPostsById = (apiDomain: string, id: string) => httpRequest.get(`${apiDomain}/${prefix}/${id}`)
export const getPosts = (apiDomain: string, pq: PageRequest) => httpRequest.get(`${apiDomain}/${prefix}`, pq)
export const getLatestPosts = (apiDomain: string) => httpRequest.get(`${apiDomain}/${prefix}/latest`)
export const likePost = (apiDomain: string, id: string) => httpRequest.post(`${apiDomain}/${prefix}/${id}/likes`, null)

