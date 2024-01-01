export interface IPost {
    id: string;
    coverImg: string;
    title: string;
    summary: string;
    categories: string[];
    tags: string[];
    createTime: number;
    updateTime: number;
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