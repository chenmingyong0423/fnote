import httpRequest from "./http";

export interface IPost {
    sug: string;
    author: string;
    title: string;
    summary: string;
    cover_img: string;
    category: string;
    tags: string[];
    like_count: number;
    comment_count: number;
    visit_count: number;
    priority: number;
    create_time: number;
}
const prefix = "/posts"

const getLatestPosts = () => {
    return httpRequest.get(prefix + "/latest")
};

export type PageRequest = {
    pageNo: number;
    pageSize: number;
    sortField?: string;
    sortOrder?: string;
    search?: string;
    category?: string;
    tags?: string[];
}

const getPosts = (pq: PageRequest) => {
    return httpRequest.get(prefix + "", pq)
};


export interface IPostDetail {
    sug: string;
    author: string;
    title: string;
    summary: string;
    cover_img: string;
    category: string;
    tags: string[];
    like_count: number;
    comment_count: number;
    visit_count: number;
    priority: number;
    create_time: number;
    content: string;
    meta_description: string;
    meta_keywords: string;
    word_count: number;
    update_time: number;
    is_liked: boolean;
}


const getPostsById = (sug: string) => {
    return httpRequest.get(prefix + "/" + sug)
};

export {
    getLatestPosts,
    getPosts,
    getPostsById
}