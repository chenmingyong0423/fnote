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

const getFriends = () => {
    return httpRequest.get(prefix + "")
};

export type FriendReq = {
    name: string;
    url: string;
    logo: string;
    description: string;
    email?: string;
}

const applyForFriend = (req: FriendReq) => {
    return httpRequest.post(prefix + "", req)
};

export {
    getFriends,
    applyForFriend
}