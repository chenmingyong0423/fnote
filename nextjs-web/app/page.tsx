import LatestComments from "@/src/components/LatestComments";
import FeaturedCarousel from "../src/components/FeaturedCarousel";
import LatestArticles from "../src/components/LatestArticles";
import SiteOwnerCard from "../src/components/SiteOwnerCard";

const featuredArticles = [
	{
		id: 1,
		title: "深入理解 Next.js 13 App Router",
		cover: "/file.svg",
		description: "全新架构，拥抱未来的 React Web 开发方式。",
		link: "/articles/1",
	},
	{
		id: 2,
		title: "TypeScript 最佳实践",
		cover: "/window.svg",
		description: "让你的前端项目更健壮、更易维护。",
		link: "/articles/2",
	},
	{
		id: 3,
		title: "Tailwind CSS 高级用法",
		cover: "/globe.svg",
		description: "极致灵活的原子化 CSS 框架。",
		link: "/articles/3",
	},
];

const latestArticles = [
	{
		id: 4,
		title: "React 18 新特性全解读",
		summary: "Concurrent 模式、自动批处理等新特性详解。",
		cover: "/next.svg",
		link: "/articles/4",
		views: 1234,
		likes: 88,
		comments: 12,
		date: "2024-06-19",
	},
	{
		id: 5,
		title: "前端性能优化实战",
		summary: "从构建到运行时的全链路优化方案。",
		cover: "/vercel.svg",
		link: "/articles/5",
		views: 2345,
		likes: 99,
		comments: 22,
		date: "2024-06-18",
	},
	{
		id: 6,
		title: "Node.js 服务端开发入门",
		summary: "快速搭建高性能 API 服务。",
		cover: "/file.svg",
		link: "/articles/6",
		views: 3456,
		likes: 77,
		comments: 33,
		date: "2024-06-17",
	},
];

const latestComments = [
	{
		id: 1,
		user: "Alice",
		avatar: "/logo.png",
		content: "文章写得很棒，受益匪浅！",
		article: { title: "深入理解 Next.js 13 App Router", link: "/articles/1" },
	},
	{
		id: 2,
		user: "Bob",
		avatar: "/logo.png",
		content: "期待更多 TypeScript 相关内容。",
		article: { title: "TypeScript 最佳实践", link: "/articles/2" },
	},
];

export default function Home() {
	return (
		<div className="w-4/5 mx-auto p-8 flex flex-col md:flex-row gap-8">
			{/* 左侧主内容区 */}
			<div className="flex-[7] flex flex-col gap-8 min-w-0">
				<FeaturedCarousel articles={featuredArticles} />
				<LatestArticles articles={latestArticles} />
			</div>
			{/* 右侧侧边栏 */}
			<div className="flex-[3] flex flex-col gap-8 shrink-0 min-w-0">
				<SiteOwnerCard />
				{/* 这里可以继续插入最新评论等其他组件 */}
        {/* 最新评论 */}
				<section className="bg-white rounded-lg shadow p-6 ">
					<h2 className="text-xl font-bold mb-4 border-b border-gray-200 pb-2">
						最新评论
					</h2>
					<LatestComments comments={latestComments} />
				</section>
			</div>
		</div>
	);
}
