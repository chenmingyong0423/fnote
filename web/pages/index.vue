<template>
  <div class="flex m-auto h-auto slide-up">
    <div class="w-69% mr-1% flex flex-col lt-md:w-100%">
      <Carousel :carousel="carousel" class="mb-5" />
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
import { getLatestPosts, getPostsById, type IPostDetail } from "~/api/post";
import type { IPost } from "~/api/post";
import type { IResponse, IListData } from "~/api/http";
import { useConfigStore } from "~/store/config";
import Carousel from "~/components/post/carousel/Carousel.vue";
import {
  type CarouselVO,
  GetCarousel,
  getInitializationStatus,
  type InitializationStatusVO,
} from "~/api/config";

const { data: posts } = await useAsyncData<IPost[]>(`latest-post`, async () => {
  let postRes: any = await getLatestPosts();
  let res: IResponse<IListData<IPost>> = postRes.data.value;
  if (res) {
    return res.data?.list || [];
  }
  return [];
});

const { data: carousel } = await useAsyncData(`get-carousel`, async () => {
  const res: any = await GetCarousel();
  const apiRes: IResponse<IListData<CarouselVO>> = res.data.value;
  if (apiRes) {
    if (apiRes.code === 0) {
      return apiRes.data?.list || [];
    }
  }
  return [];
});
</script>

<style scoped>
@keyframes slideUp {
  0% {
    transform: translateY(+100%);
  }
  100% {
    transform: translateY(0);
  }
}

.slide-up {
  animation: slideUp 0.5s ease;
  animation-iteration-count: 1;
}
</style>
