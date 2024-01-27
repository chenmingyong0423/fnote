<template>
  <div class="flex">
    <div class="w-69% mr-1% flex flex-col lt-md:w-100%">
      <div class="flex flex-col">
        <div class="mb-10 flex gap-x-5 dark:text-dtc">
          <NuxtLink
            class="p-2 cursor-pointer"
            :class="
              filter === 'latest'
                ? 'custom_bottom_border_1E80FF'
                : 'hover:custom_bottom_border_1E80FF'
            "
            :to="path + generateQuery('latest')"
            >最新
          </NuxtLink>
          <NuxtLink
            class="p-2 cursor-pointer"
            :class="
              filter === 'oldest'
                ? 'custom_bottom_border_1E80FF'
                : 'hover:custom_bottom_border_1E80FF'
            "
            :to="path + generateQuery('oldest')"
            >最早
          </NuxtLink>
          <NuxtLink
            class="p-2 cursor-pointer"
            :class="
              filter === 'likes'
                ? 'custom_bottom_border_1E80FF'
                : 'hover:custom_bottom_border_1E80FF'
            "
            :to="path + generateQuery('likes')"
            >点赞最多
          </NuxtLink>
        </div>
        <PostListItem :posts="posts" class="lt-md:hidden"></PostListItem>
        <PostListSquareItem
          :posts="posts"
          class="md:hidden"
        ></PostListSquareItem>
        <Pagination
          :currentPage="req.pageNo"
          :total="totalPosts"
          :perPageCount="req.pageSize"
          :route="`/tags/${routeParam}/page/`"
          :filterCond="filter"
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
import type { PageRequest } from "~/api/post";
import type { IPost } from "~/api/post";
import type { IResponse, IPageData } from "~/api/http";
import type { ICategoryName } from "~/api/category";
import { getTagByRoute } from "~/api/tag";
import { useConfigStore } from "~/store/config";

const route = useRoute();
const pageNo: number = +route.params.pageNo;
let pageSize: number = Number(route.query.pageSize) || 5;
const routeParam: string = String(route.params.name);
const path: string = route.path.substring(0, route.path.length - 1) + "1";
const generateQuery = (filter: string) => {
  if (pageSize != 5) {
    return "?pageSize=" + pageSize + "&filter=" + filter;
  }
  return "?filter=" + filter;
};

let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: pageNo,
  pageSize: pageSize,
  sortField: "create_time",
  sortOrder: "DESC",
} as PageRequest);
const filter = ref<string>("latest");
if (route.query.filter && route.query.filter !== "") {
  if (route.query.filter !== "latest") {
    filter.value = String(route.query.filter);
    console.log(filter.value);
    switch (filter.value) {
      case "oldest":
        console.log(12344);
        req.value.sortField = "create_time";
        req.value.sortOrder = "ASC";
        break;
      case "likes":
        req.value.sortField = "like_count";
        req.value.sortOrder = "DESC";
        break;
      default:
        req.value.sortField = "create_time";
        req.value.sortOrder = "DESC";
        break;
    }
  }
}
const totalPosts = ref<Number>(0);
const postInfos = async () => {
  try {
    if (!req.value.tags || req.value.tags.length == 0) {
      let categoryRes: any = await getTagByRoute(routeParam);
      let res: IResponse<ICategoryName> = categoryRes.data.value;
      req.value.tags = [res.data?.name || ""];
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

const configStore = useConfigStore();
useHead({
  title: `${routeParam} - ${configStore.seo_meta_config.title}`,
  meta: [
    { name: "description", content: `${routeParam}文章列表` },
    { name: "keywords", content: configStore.seo_meta_config.keywords },
    { name: "author", content: configStore.seo_meta_config.author },
    { name: "robots", content: configStore.seo_meta_config.robots },
  ],
  link: [
    { rel: "icon", type: "image/x-icon", href: configStore.website_info.icon },
  ],
});
useSeoMeta({
  ogTitle: `${routeParam} - ${configStore.seo_meta_config.og_title}`,
  ogDescription: `${routeParam}文章列表`,
  ogImage: configStore.seo_meta_config.og_image,
  twitterCard: "summary",
});

// 创建一个计算属性来追踪 query 对象
const routeQuery = computed(() => route.query);

watch(
  () => routeQuery,
  async (newQuery, oldQuery) => {
    const p: number = Number(route.query.pageSize) || -1;
    if (p != req.value.pageSize && p != -1) {
      req.value.pageSize = p;
      pageSize = p;
      await postInfos();
    }
    if (
      newQuery.value.filter &&
      newQuery.value.filter !== "" &&
      newQuery.value.filter !== filter.value
    ) {
      filter.value = String(newQuery.value.filter);
      switch (filter.value) {
        case "oldest":
          req.value.sortField = "create_time";
          req.value.sortOrder = "ASC";
          break;
        case "likes":
          req.value.sortField = "like_count";
          req.value.sortOrder = "DESC";
          break;
        default:
          req.value.sortField = "create_time";
          req.value.sortOrder = "DESC";
          break;
      }
      await postInfos();
    }
  },
  { deep: true },
);
</script>
