import httpRequest from "../utils/http";

export interface VisitLogRequest {
  url: string;
  ip?: string;
  userAgent?: string;
  origin?: string;
  referer?: string;
}

export const collectVisitLog = (req: VisitLogRequest) =>
  httpRequest.post("/logs", req);

export interface WebsiteCountStats {
  post_count: number;
  category_count: number;
  tag_count: number;
  comment_count: number;
  like_count: number;
  website_view_count: number;
}

export const getWebsiteCountStats = () => httpRequest.get("/stats");
