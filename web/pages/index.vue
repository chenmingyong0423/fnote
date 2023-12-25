<template>
  <div class="flex m-auto h-auto">
    <div class="w-69% mr-1% flex flex-col">
      <Notice></Notice>
      <PostListSquareItem :posts="posts"></PostListSquareItem>
    </div>
    <div class="flex flex-col w-30%">
      <Profile class="mb-5"></Profile>
      <IndexComment></IndexComment>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {getLatestPosts} from "~/api/post"
import type {IPost} from "~/api/post"
import type {IResponse, IListData} from "~/api/http";
let posts = ref<IPost[]>([]);

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