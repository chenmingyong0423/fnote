import React from "react";
import { getCategories } from "@/src/api/category";
import { TagsOutlined, AppstoreOutlined, ProductOutlined, BookOutlined } from "@ant-design/icons";

export default async function NavigationPage() {
    // 获取所有分类
    const categories = await getCategories();
    // 提取所有标签（假设每个分类有 tags 字段，或需单独接口获取）
    // 这里先 mock
    const tags = [
        "Next.js", "TypeScript", "React", "Node.js", "CSS", "Tailwind", "API", "SSR", "Zustand", "Ant Design"
    ];

    return (
        <div className="w-full max-w-7xl mx-auto py-8 flex flex-col gap-12 bg-white p-4 rounded-xl shadow-sm">
            {/* 分类区域 */}
            <section>
                <h2 className="text-2xl font-bold mb-6 flex items-center gap-2"><AppstoreOutlined /> 分类导航</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-6">
                    {categories.map((cat) => (
                        <a
                            key={cat.route}
                            href={`/categories${cat.route}`}
                            className="group block bg-white rounded-xl p-5 transition-transform duration-200 hover:-translate-y-2 hover:shadow-lg border border-gray-200 h-full flex flex-col"
                        >
                            <div className="flex flex-col items-start mb-2">
                                <ProductOutlined className="mb-1 text-3xl" style={{ color: '#4b5563' }} />
                                <span className="text-3xl font-bold text-black text-left">{cat.name}</span>
                            </div>
                            <div className="text-xs text-black truncate mb-2 text-left">{cat.description}</div>
                            <div className="flex items-center gap-1 text-xs text-black mt-auto">
                                <BookOutlined />
                                <span>{cat.count}</span>
                            </div>
                        </a>
                    ))}
                </div>
            </section>
            {/* 标签区域 */}
            <section >
                <h2 className="text-2xl font-bold mb-6 flex items-center gap-2"><TagsOutlined /> 标签导航</h2>
                <div className="flex flex-wrap gap-x-3 gap-y-3">
                    {tags.map((tag) => (
                        <span
                            key={tag}
                            className="px-3 py-1 rounded-xl text-gray-700 font-medium cursor-pointer transition-transform duration-200 hover:-translate-y-2 border border-gray-200"
                        >
                            #{tag}
                        </span>
                    ))}
                </div>
            </section>
        </div>
    );
}
