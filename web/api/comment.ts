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


export const getLatestComments = (apiDomain: string) => httpRequest.get(`${apiDomain}/${prefix}/latest`)

export const getComments = (apiDomain: string, id: string) => httpRequest.get(`${apiDomain}/${prefix}/id/${id}`)

export const submitComment = (apiDomain: string, comment: ICommentRequest) => httpRequest.post(`${apiDomain}/${prefix}`, comment)

export const submitCommentReply = (apiDomain: string, commentId: string, comment: ICommentReplyRequest) => httpRequest.post(`${apiDomain}/${prefix}/${commentId}/replies`, comment)
