import { getCategories } from "@/src/api/category";
import { getTags } from "@/src/api/tags";
import type { Metadata } from "next";
import {getCommonConfig} from "@/src/api/config";
import NavigationContent from "@/src/components/NavigationContent";

export async function generateMetadata(): Promise<Metadata> {
    const config = await getCommonConfig();
    const categories = await getCategories();
    const tags = await getTags();
    const catNames = categories.slice(0, 10).map((c) => c.name).join("、") + (categories.length > 10 ? "等" : "");
    const tagNames = tags.slice(0, 10).map((t) => t.name).join("、") + (tags.length > 10 ? "等" : "");
    const description = `本站导航页，包含${categories.length}个分类（${catNames}），${tags.length}个标签（${tagNames}），快速浏览所有的文章分类与标签。`;
    return {
        title: `全部分类与标签 - ${config.seo_meta.title || config.website_meta.website_name}`,
        description,
        openGraph: {
            title: `全部分类与标签 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
            description,
            url: process.env.BASE_HOST + "/navigation",
            images: config.seo_meta.og_image ? [{ url: process.env.NEXT_PUBLIC_SERVER_HOST + config.seo_meta.og_image }] : undefined,
            siteName: config.website_meta.website_name,
            type: "website",
        },
    };
}

export default async function NavigationPage() {
    // 获取所有分类
    const categories = await getCategories();
    // 获取所有标签
    const tags = await getTags();

    return <NavigationContent categories={categories} tags={tags} />;
}
