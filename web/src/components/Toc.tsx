import React from "react";

export interface TocItem {
  id: string;
  text: string;
  level: number;
}

export function extractToc(markdown: string): TocItem[] {
  const lines = markdown.split("\n");
  const toc: TocItem[] = [];
  const headingRegex = /^(#{1,6})\s+(.+)/;
  lines.forEach((line) => {
    const match = line.match(headingRegex);
    if (match) {
      const level = match[1].length;
      const text = match[2].trim();
      const id = text.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, "-").replace(/^-+|-+$/g, "");
      toc.push({ id, text, level });
    }
  });
  return toc;
}

export const Toc: React.FC<{ toc: TocItem[] }> = ({ toc }) => {
  if (!toc.length) return null;
  return (
    <nav className="sticky top-24 max-h-[70vh] overflow-auto bg-white dark:bg-[#232426] rounded-lg p-4 shadow-sm border border-gray-100 dark:border-gray-700 text-sm">
      <div className="font-bold mb-3 text-base text-gray-700 dark:text-gray-200 flex items-center gap-2">
        <svg width="18" height="18" fill="none" viewBox="0 0 24 24">
          <path
            stroke="currentColor"
            strokeWidth="2"
            d="M6 7h12M6 12h8m-8 5h12"
            strokeLinecap="round"
          />
        </svg>
        目录
      </div>
      <ul className="space-y-1">
        {toc.map((item) => (
          <li key={item.id} style={{ marginLeft: (item.level - 1) * 16 }}>
            <a
              href={`#${item.id}`}
              className="block px-2 py-1 rounded transition-colors duration-150 hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-[#1e293b] dark:hover:text-blue-400 text-gray-700 dark:text-gray-300"
              style={{
                fontWeight: item.level === 1 ? 600 : 400,
                fontSize: item.level === 1 ? "1rem" : "0.95rem",
              }}
            >
              {item.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
};
