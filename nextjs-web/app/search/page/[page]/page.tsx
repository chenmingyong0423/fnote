"use client";
import ArticleListContainer from "@/src/components/ArticleListContainer";
import { useSearchParams, useRouter } from "next/navigation";
import { Input } from "antd";

export default function SearchPageWithPagination({ params }: { params: { page: string } }) {
  const searchParams = useSearchParams();
  const router = useRouter();
  const field = (searchParams.get("filter") as "latest" | "oldest" | "likes") || "latest";
  const page = Number(params.page || 1);
  const pageSize = Number(searchParams.get("pageSize") || 10);
  const keyword = searchParams.get("keyword") || "";

  const handleSearch = (value: string) => {
    const newParams = new URLSearchParams(searchParams.toString());
    if (value) {
      newParams.set("keyword", value);
    } else {
      newParams.delete("keyword");
    }
    // 跳转到第一页
    router.replace(`/search?pageSize=${pageSize}&filter=${field}&keyword=${encodeURIComponent(value)}`);
  };

  return (
    <div className="w-full max-w-7xl mx-auto">
      <div className="mb-6 flex md:justify-start justify-center">
        <div className="w-full md:col-span-8" style={{ maxWidth: '66.6667%' }}>
          <Input.Search
            placeholder="请输入关键词..."
            allowClear
            enterButton="搜索"
            size="large"
            defaultValue={keyword}
            onSearch={handleSearch}
          />
        </div>
      </div>
      <ArticleListContainer keyword={keyword} field={field} page={page} pageSize={pageSize} />
    </div>
  );
}
