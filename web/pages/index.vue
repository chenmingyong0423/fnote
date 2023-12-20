<template>
  <div class="flex w-90% m-auto h-auto">
    <div class="w-69% mr-1% flex flex-col">
      <Notice></Notice>
      <PostListItem :posts="posts"></PostListItem>
    </div>
    <div class="w-30%">
      <Profile></Profile>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {getLatestPosts} from "~/server/api/post"
import type {IPost} from "~/server/api/post"
import type {IResponse, IListData} from "~/server/api/http";
import {onMounted, ref} from "vue";

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