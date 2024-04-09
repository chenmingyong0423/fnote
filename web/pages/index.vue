<template>
  <div class="flex m-auto h-auto">
    <div class="w-69% mr-1% flex flex-col lt-md:w-100%">
      <Notice></Notice>
      <PostListSquareItem :posts="posts"></PostListSquareItem>
      <NuxtLink class="m-auto" to="/navigation">
        <Button
            name="查看更多"
            class="w-30 h-10 line-height-10 bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200"
        ></Button>
      </NuxtLink>
    </div>
    <div class="flex flex-col w-30% lt-md:hidden">
      <Profile class="mb-5"></Profile>
      <CommentHome></CommentHome>
    </div>
  </div>
  <div>
    <ExternalLink class="mb-5"></ExternalLink>
  </div>
</template>

<script lang="ts" setup>
import {getLatestPosts} from "~/api/post";
import type {IPost} from "~/api/post";
import type {IResponse, IListData} from "~/api/http";
import {useHomeStore} from "~/store/home";
import {useConfigStore} from "~/store/config";

let posts = ref<IPost[]>([]);

const postInfos = async () => {
  try {
    let postRes: any = await getLatestPosts();
    let res: IResponse<IListData<IPost>> = postRes.data.value;
    if (res) {
      posts.value = res.data?.list || [];
    }
  } catch (error) {
    console.log(error);
  }
};
await postInfos();
const configStore = useConfigStore();
useHead({
  title: configStore.seo_meta_config.title,
  meta: [
    {name: "description", content: configStore.seo_meta_config.description},
    {name: "keywords", content: configStore.seo_meta_config.keywords},
    {name: "author", content: configStore.seo_meta_config.author},
    {name: "robots", content: configStore.seo_meta_config.robots},
  ],
  link: [
    {rel: "icon", type: "image/x-icon", href: configStore.website_info.icon},
  ],
});
useSeoMeta({
  ogTitle: configStore.seo_meta_config.og_title,
  ogDescription: configStore.seo_meta_config.description,
  ogImage: configStore.seo_meta_config.og_image,
  twitterCard: "summary",
});
</script>
