<template>
  <div class="flex">
    <div class="w-69% mr-1% flex flex-col lt-md:w-100%">
      <div class="flex flex-col">
        <PostListItem :posts="posts" class="lt-md:hidden"></PostListItem>
        <PostListSquareItem
          :posts="posts"
          class="md:hidden"
        ></PostListSquareItem>
        <Pagination
          :currentPage="req.pageNo"
          :total="totalPosts"
          :perPageCount="req.pageSize"
          :route="`/categories/${routeParam}/page/`"
        ></Pagination>
      </div>
    </div>
    <div class="flex flex-col w-30% lt-md:hidden">
      <Profile class="mb-5"></Profile>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { getPosts } from "~/api/post";
import { getCategoryByRoute } from "~/api/category";
import type { PageRequest } from "~/api/post";
import type { IPost } from "~/api/post";
import type { IResponse, IPageData } from "~/api/http";
import type { ICategoryName } from "~/api/category";
import { useHomeStore } from "~/store/home";

const route = useRoute();
const pageSize: number = Number(route.query.pageSize) || 5;
const routeParam: string = String(route.params.name);
const homeStore = useHomeStore();

let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: 1,
  pageSize: pageSize,
  sortField: "create_time",
  sortOrder: "desc",
} as PageRequest);

const totalPosts = ref<Number>(0);

const title = ref<string>("");

const postInfos = async () => {
  try {
    if (!req.value.categories || req.value.categories.length == 0) {
      let categoryRes: any = await getCategoryByRoute(routeParam);
      let res: IResponse<ICategoryName> = categoryRes.data.value;
      title.value = res.data?.name || "";
      req.value.categories = [title.value];
    }
    const deepCopyReq = JSON.parse(JSON.stringify(req.value));
    let postRes: any = await getPosts(deepCopyReq);
    let res: IResponse<IPageData<IPost>> = postRes.data.value;
    posts.value = res.data?.list || [];
    totalPosts.value = res.data?.totalCount || totalPosts.value;
  } catch (error) {
    console.log(error);
  }
};
await postInfos();

useHead({
  title: `${title.value} - ${homeStore.seo_meta_config.title}`,
  meta: [
    { name: "description", content: `${title.value}文章列表` },
    { name: "keywords", content: homeStore.seo_meta_config.keywords },
    { name: "author", content: homeStore.seo_meta_config.author },
    { name: "robots", content: homeStore.seo_meta_config.robots },
  ],
  link: [
    { rel: "icon", type: "image/x-icon", href: homeStore.website_info.icon },
  ],
});
useSeoMeta({
  ogTitle: `${title.value} - ${homeStore.seo_meta_config.og_title}`,
  ogDescription: `${title.value}文章列表`,
  ogImage: homeStore.seo_meta_config.og_image,
  twitterCard: "summary",
});
</script>
