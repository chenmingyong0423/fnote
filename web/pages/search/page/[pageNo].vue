<template>
  <div class="flex">
    <div class="w-69% mr-1% flex flex-col lt-md:w-100%">
      <div class="flex flex-col">
        <SearchInput :keyword="keyword" class="mb-5"></SearchInput>
        <div class="mb-10 flex gap-x-5">
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
          :route="'/search/page/'"
          :extraParams="{ keyword: req.keyword }"
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
import { useConfigStore } from "~/store/config";

const route = useRoute();
const pageNo: number = +route.params.pageNo;
let pageSize: number = Number(route.query.pageSize) || 5;
let keyword = ref<string>(String(route.query.keyword));
if (keyword.value == "undefined") {
  keyword.value = "";
}
const path: string = route.path.substring(0, route.path.length - 1) + "1";
const generateQuery = (filter: string) => {
  if (pageSize != 5) {
    return (
      "?pageSize=" +
      pageSize +
      "&keyword=" +
      keyword.value +
      "&filter=" +
      filter
    );
  }
  return "?keyword=" + keyword.value + "&filter=" + filter;
};
let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: pageNo,
  pageSize: pageSize,
  sortField: "create_time",
  sortOrder: "desc",
  keyword: keyword.value,
} as PageRequest);

const totalPosts = ref<Number>(0);
const filter = ref<string>("latest");
if (route.query.filter && route.query.filter !== "") {
  if (route.query.filter !== "latest") {
    filter.value = String(route.query.filter);
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
  }
}
const postInfos = async () => {
  try {
    const deepCopyReq = JSON.parse(JSON.stringify(req.value));
    let postRes: any = await getPosts(deepCopyReq);
    let res: IResponse<IPageData<IPost>> = postRes.data.value;
    posts.value = res.data?.list || [];
    totalPosts.value = res.data?.totalCount || 0;
  } catch (error) {
    console.log(error);
  }
};
await postInfos();

// 创建一个计算属性来追踪 query 对象
const routeQuery = computed(() => route.query);

watch(
  () => routeQuery,
  async (newQuery, oldQuery) => {
    const p: number = Number(newQuery.value.pageSize) || -1;
    if (p != req.value.pageSize && p != -1) {
      req.value.pageSize = p;
      pageSize = p;
      await postInfos();
    } else if (newQuery.value.keyword !== keyword.value) {
      req.value.keyword = String(route.query.keyword);
      keyword.value = String(route.query.keyword);
      await postInfos();
      seo();
    } else if (
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

const configStore = useConfigStore();
const seo = () => {
  useHead({
    title: `${keyword.value} - 搜索 - ${configStore.seo_meta_config.title}`,
    meta: [
      { name: "description", content: `${keyword.value} 搜索结果` },
      { name: "keywords", content: configStore.seo_meta_config.keywords },
      { name: "author", content: configStore.seo_meta_config.author },
      { name: "robots", content: configStore.seo_meta_config.robots },
    ],
    link: [
      {
        rel: "icon",
        type: "image/x-icon",
        href: configStore.website_info.icon,
      },
    ],
  });
  useSeoMeta({
    ogTitle: `${keyword.value} - 搜索 - ${configStore.seo_meta_config.og_title}`,
    ogDescription: `${keyword.value} 搜索结果`,
    ogImage: configStore.seo_meta_config.og_image,
    twitterCard: "summary",
  });
};
seo();
</script>
