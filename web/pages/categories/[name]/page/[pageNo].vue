<template>
  <div class="flex flex-col">
    <PostListItem :posts="posts"></PostListItem>
    <Pagination :currentPage="req.pageNo" :total="totalPosts" :perPageCount='req.pageSize'
                :route="'/categories/' + pathName + '/page/'"></Pagination>
  </div>
</template>


<script lang="ts" setup>
import {getPosts} from "~/api/post"
import type {PageRequest} from "~/api/post"
import type {IPost} from "~/api/post"
import type {IResponse, IPageData} from "~/api/http";
import {useHomeStore} from '~/store/home';
import type {IMenu} from "~/api/category";

const homeStore = useHomeStore()
const route = useRoute()
const pathName = route.params.name
const pageNo : number = +route.params.pageNo
const pageSize :number = Number(route.query.pageSize) || 5


let name = homeStore.menuList.find((item: IMenu) => item.route == '/categories/' + pathName)?.name
let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: pageNo,
  pageSize: pageSize,
  sortField: "create_time",
  sortOrder: "desc",
  categories: [name],
} as PageRequest)

const totalPosts = ref<Number>(0)

const postInfos = async () => {
  try {
    const deepCopyReq = JSON.parse(JSON.stringify(req.value));
    let postRes: any = await getPosts(deepCopyReq)
    let res: IResponse<IPageData<IPost>> = postRes.data.value
    posts.value = res.data?.list || []
    totalPosts.value = res.data?.totalCount || totalPosts.value
  } catch (error) {
    console.log(error);
  }
};
postInfos()

// 创建一个计算属性来追踪 query 对象
const routeQuery = computed(() => route.query);

watch(() => routeQuery, (newQuery, oldQuery) => {
  const pageSize :number = Number(route.query.pageSize) || -1
  if (pageSize != req.value.pageSize && pageSize != -1){
    req.value.pageSize = pageSize
    postInfos()
  }
}, { deep: true });
</script>