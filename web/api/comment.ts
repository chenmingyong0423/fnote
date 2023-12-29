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
export const getComments = (id: string) => httpRequest.get(`${prefix}/id/${id}`)

export interface ICommentRequest {
    postId: string;
    username: string;
    email: string;
    website?: string;
    content: string;
}

export const submitComment = (comment: ICommentRequest) => httpRequest.post(`${prefix}`, comment)

export interface ICommentReplyRequest {
    postId: string;
    replyToId?: string;
    username: string;
    email: string;
    website?: string;
    content: string;
}

export const submitCommentReply = (commentId: string, comment: ICommentReplyRequest) => httpRequest.post(`${prefix}/${commentId}/replies`, comment)
