import httpRequest from "./http";

export interface IFriend {
    name: string;
    url: string;
    logo: string;
    description: string;
    website_live_time: string;
    priority: string;
}
const prefix = "/friends"

export const getFriends = () => httpRequest.get(prefix + "")


export type FriendReq = {
    name: string;
    url: string;
    logo: string;
    description: string;
    email?: string;
}

export const applyForFriend = (req: FriendReq) => httpRequest.post(prefix + "", req)


