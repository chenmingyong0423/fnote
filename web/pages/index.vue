<template>
  <div class="flex m-auto h-auto">
    <div class="w-69% mr-1% flex flex-col">
      <Notice></Notice>
      <PostListItem :posts="posts"></PostListItem>
<!--      <Pagination :currentPage="req.pageNo" :total="totalPosts" :perPageCount='req.pageSize' @pageChanged="pageChanged"-->
<!--                  @perPageChanged="perPageChanged"></Pagination>-->
    </div>
    <div class="flex flex-col w-30%">
      <Profile class="mb-5"></Profile>
      <IndexComment></IndexComment>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {getPosts} from "~/server/api/post"
import type {PageRequest} from "~/server/api/post"
import type {IPost} from "~/server/api/post"
import type {IResponse, IPageData} from "~/server/api/http";


let posts = ref<IPost[]>([]);
let req = reactive<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: "create_time",
  sortOrder: "desc",
} as PageRequest)


// const totalPosts = ref<Number>(0)

const pageChanged = (page: number) => {
  req.pageNo = page
  postInfos()
}

const perPageChanged = (itemsPerPage: number) => {
  req.pageNo = 1
  req.pageSize = itemsPerPage
  postInfos()
}

const postInfos = async () => {
  try {
    const deepCopyReq =  JSON.parse(JSON.stringify(req));
    let postRes: any = await getPosts(deepCopyReq)
    let res: IResponse<IPageData<IPost>> = postRes.data.value
    posts.value = res.data?.list || []
    // totalPosts.value = res.data?.totalCount || totalPosts.value
  } catch (error) {
    console.log(error);
  }
};
postInfos()
</script>