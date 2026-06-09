"use client";
import { useEffect } from "react";
import { logVisit } from "@/src/api/logs";

export default function LogVisitClient() {
  useEffect(() => {
    const url = window.location.href;
    const user_agent = navigator.userAgent;
    const origin = window.location.origin;
    const referer = document.referrer;
    void logVisit({ url, user_agent, origin, referer }).catch(() => {
      // Visit logging is best-effort and should not break page rendering.
    });
  }, []);
  return null;
}
