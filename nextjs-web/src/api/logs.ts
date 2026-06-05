import { request } from "../utils/http";
import type { Response } from "./types";

export interface VisitLogReq {
  url: string;
  ip?: string;
  user_agent?: string;
  origin?: string;
  referer?: string;
}

export async function logVisit(data: VisitLogReq): Promise<Response<null>> {
  return request<Response<null>>("/api/logs", {
    method: "POST",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json" },
  });
}
