import httpRequest from "~/api/http";

export interface VisitLogRequest {
    url: string;
    ip: string;
    userAgent: string;
    origin: string;
    referer: string;
}

export const collectVisitLog = (req: VisitLogRequest) => httpRequest.post("/logs", req)