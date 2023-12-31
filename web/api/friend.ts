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

export const getFriends = (apiDomain: string) => httpRequest.get(`${apiDomain}/${prefix}`)
export const applyForFriend = (apiDomain: string, req: FriendReq) => httpRequest.post(`${apiDomain}/${prefix}`, req)


