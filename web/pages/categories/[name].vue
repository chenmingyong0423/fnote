<template>
  <div class="flex flex-col">
    <PostListItem :posts="posts"></PostListItem>
    <Pagination :currentPage="req.pageNo" :total="totalPosts" :perPageCount='req.pageSize'
                :route="path +'/page/'"></Pagination>
  </div>
</template>


<script lang="ts" setup>
import {getPosts} from "~/server/api/post"
import type {PageRequest} from "~/server/api/post"
import type {IPost} from "~/server/api/post"
import type {IResponse, IPageData} from "~/server/api/http";
import {useHomeStore} from '~/store/home';
import type {IMenu} from "~/server/api/category";

const homeStore = useHomeStore()
const route = useRoute()
const path = route.path
let name = homeStore.menuList.find((item: IMenu) => item.route == path)?.name
console.log(123)
let posts = ref<IPost[]>([]);
let req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: "create_time",
  sortOrder: "desc",
  category: name,
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
</script>