"use client";

import { NotificationOutlined } from "@ant-design/icons";
import type { NoticeConfigVO } from "@/src/api/config";

interface HomeNoticeProps {
  notice: NoticeConfigVO;
  hasError?: boolean;
}

export default function HomeNotice({ notice, hasError = false }: HomeNoticeProps) {
  const title = notice.title?.trim();
  const content = notice.content?.trim();

  if (hasError || !notice.enabled || (!title && !content)) {
    return null;
  }

  return (
    <section className="glass-surface rounded-lg px-4 py-3 md:px-5 md:py-4">
      <div className="flex items-start gap-3">
        <span className="mt-0.5 inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border border-sky-200 bg-sky-50/80 text-sky-700 dark:border-sky-900/70 dark:bg-sky-950/70 dark:text-sky-200">
          <NotificationOutlined />
        </span>
        <div className="min-w-0 flex-1">
          {title && (
            <h2 className="truncate text-sm font-semibold leading-6 text-gray-950 dark:text-gray-100 md:text-base">
              {title}
            </h2>
          )}
          {content && (
            <p className="line-clamp-2 text-sm leading-6 text-gray-600 dark:text-gray-300">
              {content}
            </p>
          )}
        </div>
      </div>
    </section>
  );
}
