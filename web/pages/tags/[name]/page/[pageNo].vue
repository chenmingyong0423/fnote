<template>
  <div class="flex">
    <div class="w-69% mr-1% flex flex-col">
      <div class="flex flex-col">
        <PostListItem :posts="posts"></PostListItem>
        <Pagination :currentPage="req.pageNo" :total="totalPosts" :perPageCount='req.pageSize'
                    :route="`/tags/${routeParam}/page/`"></Pagination>
      </div>
    </div>
    <div class="flex flex-col w-30%">
      <Profile class="mb-5"></Profile>
    </div>
  </div>
</template>


<script lang="ts" setup>
import {getPosts} from "~/api/post"
import type {PageRequest} from "~/api/post"
import type {IPost} from "~/api/post"
import type {IResponse, IPageData} from "~/api/http";
import type {ICategoryName} from "~/api/category";
import {getTagByRoute} from "~/api/tag";
import {useHomeStore} from "~/store/home";

const route = useRoute()
const pageNo: number = +route.params.pageNo
const pageSize: number = Number(route.query.pageSize) || 5
const routeParam: string = String(route.params.name)


let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: pageNo,
  pageSize: pageSize,
  sortField: "create_time",
  sortOrder: "desc",
} as PageRequest)

const totalPosts = ref<Number>(0)
const homeStore = useHomeStore()
const apiBaseUrl = homeStore.apiBaseUrl;
const postInfos = async () => {
  try {
    if (!req.value.tags || req.value.tags.length == 0) {
      let categoryRes: any = await getTagByRoute(apiBaseUrl, routeParam)
      let res: IResponse<ICategoryName> = categoryRes.data.value
      req.value.tags = [res.data?.name || ""]
    }
    const deepCopyReq = JSON.parse(JSON.stringify(req.value));
    let postRes: any = await getPosts(apiBaseUrl, deepCopyReq)
    let res: IResponse<IPageData<IPost>> = postRes.data.value
    posts.value = res.data?.list || []
    totalPosts.value = res.data?.totalCount || totalPosts.value
  } catch (error) {
    console.log(error);
  }
};

await postInfos()

useHead({
  title: `${routeParam} - ${homeStore.seo_meta_config.title}`,
  meta: [
    {name: 'description', content: `${routeParam}文章列表`},
    {name: 'keywords', content: homeStore.seo_meta_config.keywords},
    {name: 'author', 'content': homeStore.seo_meta_config.author},
    {name: 'robots', 'content': homeStore.seo_meta_config.robots},
  ],
  link: [
    {rel: 'icon', type: 'image/x-icon', href: homeStore.master_info.website_icon},
  ]
})
useSeoMeta({
  ogTitle: `${routeParam} - ${homeStore.seo_meta_config.title}`,
  ogDescription: `${routeParam}文章列表`,
  ogImage: '',
  twitterCard: 'summary'
})

// 创建一个计算属性来追踪 query 对象
const routeQuery = computed(() => route.query);

watch(() => routeQuery, (newQuery, oldQuery) => {
  const pageSize: number = Number(route.query.pageSize) || -1
  if (pageSize != req.value.pageSize && pageSize != -1) {
    req.value.pageSize = pageSize
    postInfos()
  }
}, {deep: true});
</script>