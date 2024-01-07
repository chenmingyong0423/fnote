import httpRequest from "./http";

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
};

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

const prefix = "/posts";
export const getPostsById = (id: string) => httpRequest.get(`${prefix}/${id}`);
export const getPosts = (pq: PageRequest) => httpRequest.get(`${prefix}`, pq);
export const getLatestPosts = () => httpRequest.get(`${prefix}/latest`);
export const likePost = (id: string) =>
  httpRequest.post(`${prefix}/${id}/likes`, null);
