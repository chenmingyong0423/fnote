import httpRequest from "./http";

export interface ILatestComment {
    post_id: string;
    post_title: string;
    picture: string;
    name: string;
    content: string;
    create_time: number;
}

const prefix = "/comments"

export const getLatestComments = () => httpRequest.get(prefix + "/latest")


export interface IReply {
    id: string;
    comment_id: string;
    content: string;
    name: string;
    reply_to_id: string;
    reply_to: string;
    reply_time: number;
    replied_content?: string;
}

export interface IComment {
    id: string;
    content: string;
    username: string;
    comment_time: number;
    replies: IReply[];
}

export const getComments = (sug: string) => httpRequest.get(prefix + "/sug/" + sug)


