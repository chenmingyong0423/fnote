import httpRequest from "./http";

export interface ILatestComment {
    post_id: string;
    post_title: string;
    picture: string;
    name: string;
    content: string;
    create_time: number;
}

export interface IReply {
    id: string;
    comment_id: string;
    content: string;
    name: string;
    email: string;
    website?: string;
    reply_to_id: string;
    reply_to: string;
    reply_time: number;
    replied_content?: string;
}

export interface IComment {
    id: string;
    content: string;
    username: string;
    email: string;
    website?: string;
    comment_time: number;
    replies: IReply[];
}

export interface ICommentRequest {
    postId: string;
    username: string;
    email: string;
    website?: string;
    content: string;
}


export interface ICommentReplyRequest {
    postId: string;
    replyToId?: string;
    username: string;
    email: string;
    website?: string;
    content: string;
}

const prefix = "comments"


export const getLatestComments = (apiBaseUrl: string) => httpRequest.get(`${apiBaseUrl}/${prefix}/latest`)

export const getComments = (apiBaseUrl: string, id: string) => httpRequest.get(`${apiBaseUrl}/${prefix}/id/${id}`)

export const submitComment = (apiBaseUrl: string, comment: ICommentRequest) => httpRequest.post(`${apiBaseUrl}/${prefix}`, comment)

export const submitCommentReply = (apiBaseUrl: string, commentId: string, comment: ICommentReplyRequest) => httpRequest.post(`${apiBaseUrl}/${prefix}/${commentId}/replies`, comment)
