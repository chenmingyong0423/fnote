import { request } from "../utils/http";
import type { Response } from "./types";

export interface CheckInitializationResp {
  initStatus: boolean;
}

export async function checkInitialization(): Promise<Response<CheckInitializationResp>> {
  return request<Response<CheckInitializationResp>>("/api/configs/check-initialization");
}
