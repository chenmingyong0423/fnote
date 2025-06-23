import LatestComments from "@/src/components/LatestComments";
import FeaturedCarousel from "../src/components/FeaturedCarousel";
import LatestPosts from "../src/components/LatestPosts";
import SiteOwnerCard from "../src/components/SiteOwnerCard";
import { getLatestPosts } from "../src/api/posts";
import { getCarouselList } from "../src/api/carousel";
import { getLatestComments } from "../src/api/comments";

export default async function Home() {
	// 动态获取最新文章
	const latestArticles = await getLatestPosts();
	// 动态获取轮播图
	const carouselItems = await getCarouselList();
	// 动态获取最新评论
	const latestComments = await getLatestComments();

	return (
		<div className="w-full max-w-7xl mx-auto py-4 grid grid-cols-1 md:grid-cols-12 gap-8">
			{/* 左侧主内容区 8/12 */}
			<div className="md:col-span-8 flex flex-col gap-8 min-w-0">
				{/* Banner 区 */}
				<section>
					<FeaturedCarousel items={carouselItems} />
				</section>
				{/* 内容卡片区 2列 */}
				<section>
					<LatestPosts articles={latestArticles} />
				</section>
			</div>
			{/* 右侧信息区 4/12 */}
			<div className="md:col-span-4 flex flex-col gap-8 min-w-0">
				{/* 用户卡片 */}
				<SiteOwnerCard />
				{/* 最新评论卡片 */}
				<LatestComments comments={latestComments} />
			</div>
		</div>
	);
}
