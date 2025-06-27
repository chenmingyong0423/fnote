"use client";
import { useEffect } from "react";
import { logVisit } from "@/src/api/logs";

export default function LogVisitClient() {
  useEffect(() => {
    // 1. 检查初始化状态
    import("@/src/api/checkInitialization").then(({ checkInitialization }) => {
      checkInitialization().then((res) => {
        if (res.code === 0 && res.data && res.data.initStatus === false) {
          // 跳转到 ADMIN_HOST
          const adminHost = process.env.NEXT_PUBLIC_ADMIN_HOST || process.env.ADMIN_HOST;
          if (adminHost) {
            window.location.href = adminHost;
            return;
          }
        }
        // 2. 正常访问日志埋点
        const url = window.location.href;
        const user_agent = navigator.userAgent;
        const origin = window.location.origin;
        const referer = document.referrer;
        logVisit({ url, user_agent, origin, referer });
      });
    });
  }, []);
  return null;
}
