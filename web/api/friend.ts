import httpRequest from "./http";

export interface IFriend {
  name: string;
  url: string;
  logo: string;
  description: string;
}

export type FriendReq = {
  name: string;
  url: string;
  logo: string;
  description: string;
  email?: string;
};

const prefix = "/friends";

export const getFriends = () => httpRequest.get(`${prefix}`);
export const applyForFriend = (req: FriendReq) =>
  httpRequest.post(`${prefix}`, req);

export interface FriendIntroductionVO {
  introduction: string;
}

export const GetFriendIntroduction = () => httpRequest.get(`${prefix}/summary`);
