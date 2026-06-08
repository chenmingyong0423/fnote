"use client";

import { AppstoreOutlined, BookOutlined, ProductOutlined, TagsOutlined } from "@ant-design/icons";
import { Empty } from "antd";
import type { CategoryWithCountVO } from "@/src/api/category";
import type { TagVO } from "@/src/api/tags";

type NavigationContentProps = {
  categories: CategoryWithCountVO[];
  tags: TagVO[];
};

export default function NavigationContent({ categories, tags }: NavigationContentProps) {
  const hasCategories = categories.length > 0;
  const hasTags = tags.length > 0;

  return (
    <div className="w-full max-w-7xl mx-auto md:px-4 py-6 md:py-8 flex flex-col gap-8 md:gap-12 bg-white p-4 md:p-6 rounded-xl shadow-sm dark:bg-[#141414] dark:text-gray-300 dark:border dark:border-[#303030]">
      <section>
        <h2 className="text-xl md:text-2xl font-bold mb-4 md:mb-6 flex items-center gap-2 dark:text-gray-300">
          <AppstoreOutlined /> 分类导航
        </h2>
        {hasCategories ? (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3 md:gap-6 dark:text-gray-300">
            {categories.map((cat) => (
              <a
                key={cat.route}
                href={`/categories/${cat.route}`}
                className="group block bg-white rounded-lg md:rounded-xl p-3 md:p-5 transition-transform duration-200 hover:-translate-y-1 md:hover:-translate-y-2 hover:shadow-lg border border-gray-200 h-full flex flex-col delay-500 dark:bg-[#232426] dark:border-gray-700 dark:hover:border-gray-600"
              >
                <div className="flex flex-col items-start mb-2 dark:text-gray-300">
                  <ProductOutlined className="mb-1 text-lg md:text-3xl" />
                  <span className="text-sm md:text-3xl font-bold text-left line-clamp-2">{cat.name}</span>
                </div>
                <div className="text-xs truncate mb-2 text-left dark:text-gray-400 line-clamp-2">{cat.description}</div>
                <div className="flex items-center gap-1 text-xs mt-auto">
                  <BookOutlined />
                  <span className="dark:text-gray-400">{cat.count}</span>
                </div>
              </a>
            ))}
          </div>
        ) : (
          <div className="py-6 flex justify-center">
            <Empty description="暂无分类" />
          </div>
        )}
      </section>

      <section>
        <h2 className="text-xl md:text-2xl font-bold mb-4 md:mb-6 flex items-center gap-2">
          <TagsOutlined /> 标签导航
        </h2>
        {hasTags ? (
          <div className="flex flex-wrap gap-2 md:gap-3">
            {tags.map((tag) => (
              <a
                key={tag.route}
                href={`/tags/${tag.route}`}
                className="px-2 md:px-3 py-1 rounded-full text-xs md:text-sm font-medium cursor-pointer transition-transform duration-200 hover:-translate-y-1 md:hover:-translate-y-2 border border-gray-200 delay-50 dark:border-gray-600 dark:text-gray-300 dark:hover:border-gray-500"
              >
                #{tag.name}
              </a>
            ))}
          </div>
        ) : (
          <div className="py-6 flex justify-center">
            <Empty description="暂无标签" />
          </div>
        )}
      </section>
    </div>
  );
}
