"use client";

import {
  AppstoreOutlined,
  BookOutlined,
  ProductOutlined,
  TagsOutlined,
} from "@ant-design/icons";
import { Empty } from "antd";
import type { CategoryWithCountVO } from "@/src/api/category";
import type { TagVO } from "@/src/api/tags";

type NavigationContentProps = {
  categories: CategoryWithCountVO[];
  tags: TagVO[];
  hasCategoryError?: boolean;
  hasTagError?: boolean;
};

export default function NavigationContent({
  categories,
  tags,
  hasCategoryError = false,
  hasTagError = false,
}: NavigationContentProps) {
  const hasCategories = categories.length > 0;
  const hasTags = tags.length > 0;

  return (
    <div className="w-full max-w-7xl mx-auto px-4 md:px-6 py-5 md:py-8 flex flex-col gap-7 md:gap-12 text-gray-900 dark:text-gray-100">
      <h1 className="text-xl font-bold text-gray-900 dark:text-gray-100 md:text-2xl">
        全部分类与标签
      </h1>
      <section>
        <h2 className="text-xl md:text-2xl font-bold mb-4 md:mb-6 flex items-center gap-2 text-gray-900 dark:text-gray-100">
          <AppstoreOutlined /> 分类导航
        </h2>
        {hasCategories ? (
          <div className="grid grid-cols-1 min-[420px]:grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3 md:gap-6">
            {categories.map((cat) => (
              <a
                key={cat.route}
                href={`/categories/${cat.route}`}
                className="glass-card group block rounded-lg p-3 md:p-5 transition-transform duration-200 md:hover:-translate-y-2 h-full flex flex-col delay-500"
              >
                <div className="flex flex-col items-start mb-2 text-gray-900 dark:text-gray-100">
                  <ProductOutlined className="mb-1 text-lg md:text-2xl" />
                  <h3 className="text-sm md:text-xl font-bold text-left line-clamp-2">
                    {cat.name}
                  </h3>
                </div>
                <div className="text-xs truncate mb-2 text-left text-gray-600 dark:text-gray-400 line-clamp-2">
                  {cat.description}
                </div>
                <div className="flex items-center gap-1 text-xs mt-auto text-gray-500 dark:text-gray-400">
                  <BookOutlined />
                  <span>{cat.count}</span>
                </div>
              </a>
            ))}
          </div>
        ) : (
          <div className="py-6 flex justify-center">
            <Empty
              description={hasCategoryError ? "网站数据暂时异常" : "暂无分类"}
            />
          </div>
        )}
      </section>

      <section>
        <h2 className="text-xl md:text-2xl font-bold mb-4 md:mb-6 flex items-center gap-2 text-gray-900 dark:text-gray-100">
          <TagsOutlined /> 标签导航
        </h2>
        {hasTags ? (
          <div className="flex flex-wrap gap-2 md:gap-3">
            {tags.map((tag) => (
              <a
                key={tag.route}
                href={`/tags/${tag.route}`}
                className="inline-flex max-w-full items-center gap-2 px-2 md:px-3 py-1 rounded-full text-xs md:text-sm font-medium cursor-pointer transition-transform duration-200 md:hover:-translate-y-2 border border-white/70 bg-white/50 text-gray-700 backdrop-blur delay-50 dark:border-white/10 dark:bg-slate-900/45 dark:text-gray-300 dark:hover:border-gray-500"
              >
                <span className="min-w-0 break-all">#{tag.name}</span>
                <span className="inline-flex shrink-0 items-center gap-1 text-[11px] text-gray-500 dark:text-gray-400">
                  <BookOutlined />
                  {tag.count}
                </span>
              </a>
            ))}
          </div>
        ) : (
          <div className="py-6 flex justify-center">
            <Empty
              description={hasTagError ? "网站数据暂时异常" : "暂无标签"}
            />
          </div>
        )}
      </section>
    </div>
  );
}
