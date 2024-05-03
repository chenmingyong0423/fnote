import httpRequest from "~/api/http";

export interface PostVisitRequest {
  post_id: string;
  stay_time: number;
  visit_at: number;
}

export const CollectPostVisit = (req: PostVisitRequest) =>
  httpRequest.post("/logs/post-visit", req);
