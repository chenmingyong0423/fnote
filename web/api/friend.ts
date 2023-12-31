import httpRequest from "./http";

export interface IFriend {
    name: string;
    url: string;
    logo: string;
    description: string;
}

const prefix = "friends"



export type FriendReq = {
    name: string;
    url: string;
    logo: string;
    description: string;
    email?: string;
}

export const getFriends = (apiBaseUrl: string) => httpRequest.get(`${apiBaseUrl}/${prefix}`)
export const applyForFriend = (apiBaseUrl: string, req: FriendReq) => httpRequest.post(`${apiBaseUrl}/${prefix}`, req)


