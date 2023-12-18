<template>
  <div class="flex">
    <div class="w-7/10">
      <PostListItem :posts="posts"></PostListItem>
    </div>
    <div class="w-3/10">
      right
    </div>
  </div>
</template>

<script lang="ts" setup>
import {getLatestPosts} from "~/server/api/post"
import type {IPost} from "~/server/api/post"
import type {IResponse, IListData} from "~/server/api/http";

const posts = ref<IPost[]>([]);
const postInfos = async () => {
  try {
    let postRes: any = await getLatestPosts()
    let res: IResponse<IListData<IPost>> = postRes.data.value
    posts.value = res.data?.list || []
  } catch (error) {
    console.log(error);
  }
};
postInfos()
</script>