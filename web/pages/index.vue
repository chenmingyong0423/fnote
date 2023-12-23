<template>
  <div class="flex m-auto h-auto">
    <div class="w-69% mr-1% flex flex-col">
      <Notice></Notice>
      <PostListItem :posts="posts"></PostListItem>
      <Pagination :currentPage="req.pageNo" :totalItems="totalPosts" :pageChanged="pageChanged" :perPageChanged="perPageChanged"></Pagination>

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
import {ref} from "vue";

const posts = ref<IPost[]>([]);
const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: "create_time",
  sortOrder: "desc",
} as PageRequest)

const totalPosts = ref<Number>(0)

const pageChanged = (page: number) => {
  req.value.pageNo = page
  postInfos()
}

const perPageChanged = (itemsPerPage: number) => {
  req.value.pageSize = itemsPerPage
  postInfos()
}

const postInfos = async () => {
  try {
    let postRes: any = await getPosts(req.value)
    let res: IResponse<IPageData<IPost>> = postRes.data.value
    posts.value = res.data?.list || []
    totalPosts.value = res.data?.totalCount || totalPosts.value
  } catch (error) {
    console.log(error);
  }
};
postInfos()
</script>