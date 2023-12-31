<template>
  <div class="flex">
    <div class="w-69% mr-1% flex flex-col lt-md:w-100%">
      <div class="flex flex-col">
        <SearchInput :keyword="keyword" class="mb-5"></SearchInput>
        <PostListItem :posts="posts"></PostListItem>
        <Pagination :currentPage="req.pageNo" :total="totalPosts" :perPageCount='req.pageSize'
                    :route="'/search/page/'" :extraParams="{keyword: req.keyword}"></Pagination>
      </div>
    </div>
    <div class="flex flex-col w-30% lt-md:hidden">
      <Profile class="mb-5"></Profile>
    </div>
  </div>

</template>


<script lang="ts" setup>
import {getPosts} from "~/api/post"
import type {PageRequest} from "~/api/post"
import type {IPost} from "~/api/post"
import type {IResponse, IPageData} from "~/api/http";
import {useHomeStore} from '~/store/home';
import type {IMenu} from "~/api/category";

const route = useRoute()
const pageNo : number = +route.params.pageNo

const pageSize: number = Number(route.query.pageSize) || 5
let keyword = ref<string>(String(route.query.keyword))
if (keyword.value == 'undefined') {
  keyword.value = ""
}

let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: pageNo,
  pageSize: pageSize,
  sortField: "create_time",
  sortOrder: "desc",
  keyword: keyword.value,
} as PageRequest)

const totalPosts = ref<Number>(0)

const postInfos = async () => {
  try {
    const deepCopyReq = JSON.parse(JSON.stringify(req.value));
    let postRes: any = await getPosts(deepCopyReq)
    let res: IResponse<IPageData<IPost>> = postRes.data.value
    posts.value = res.data?.list || []
    totalPosts.value = res.data?.totalCount || 0
  } catch (error) {
    console.log(error);
  }
};
await postInfos()

// 创建一个计算属性来追踪 query 对象
const routeQuery = computed(() => route.query);

watch(() => routeQuery, async (newQuery, oldQuery) => {
  const pageSize :number = Number(route.query.pageSize) || -1
  if (pageSize != req.value.pageSize && pageSize != -1){
    req.value.pageSize = pageSize
    keyword.value = String(route.query.keyword)
    await postInfos()
    seo()
  }
}, { deep: true });

const homeStore = useHomeStore()
const seo = ()=> {
  useHead({
    title: `${keyword.value} - 搜索 - ${homeStore.seo_meta_config.title}`,
    meta: [
      {name: 'description', content: `${keyword.value} 搜索结果`},
      { name: 'keywords', content: homeStore.seo_meta_config.keywords },
      { name: 'author', 'content': homeStore.seo_meta_config.author },
      { name: 'robots', 'content': homeStore.seo_meta_config.robots },
    ]
  })
  useSeoMeta({
    ogTitle: `${keyword.value} - 搜索 - ${homeStore.seo_meta_config.title}`,
    ogDescription: `${keyword.value} 搜索结果`,
    ogImage: '',
    twitterCard: 'summary'
  })
}
seo()
</script>