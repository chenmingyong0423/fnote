<template>
  <div class="flex">
    <div class="w-69% mr-1% flex flex-col">
      <div class="flex flex-col">
        <SearchInput :keyword="keyword" class="mb-5" @search="search"></SearchInput>
        <PostListItem :posts="posts"></PostListItem>
        <Pagination :currentPage="req.pageNo" :total="totalPosts" :perPageCount='req.pageSize'
                    :route="path +'/page/'" :extraParams="{keyword: req.keyword}"></Pagination>
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
import {useHomeStore} from '~/store/home';
import type {IMenu} from "~/api/category";

const route = useRoute()
const path = route.path
const pageSize: number = Number(route.query.pageSize) || 5
let keyword = ref<string>(String(route.query.keyword))
if (keyword.value == 'undefined') {
  keyword.value = ""
}

let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: 1,
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
postInfos()

// 创建一个计算属性来追踪 query 对象
const routeQuery = computed(() => route.query);

watch(() => routeQuery, (newQuery, oldQuery) => {
  req.value.keyword = String(route.query.keyword)
  postInfos()
}, { deep: true });
</script>