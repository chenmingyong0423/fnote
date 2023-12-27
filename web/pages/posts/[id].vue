<template>
  <div class="flex w-full">
    <div class="flex flex-col gap-y-3 w-10%  items-center mt-10">
      <div class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
        <span class="i-ph:thumbs-up w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
      </div>
      <div class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
        <span class="i-ph-chats-duotone w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
      </div>
      <div class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
        <span class="i-ph:share-fat-light w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
      </div>
    </div>
    <div class="bg-white w-58% ml-1% mr-1%">
      <!--  文章标题  -->
      <div class="text-10 font-bold text-center p-1">{{post?.title}}</div>
      <!--  文章 meta  -->
      <div class="flex items-center gap-x-2 text-4 justify-center p-1 text-gray-4">
        <div>{{post?.author}}</div>
        <div>{{ $dayjs(post?.create_time).format('YYYY-MM-DD HH:mm:ss') }}</div>
        <div>阅读 {{post?.visit_count}}</div>
      </div>
      <!--  文章内容  -->
      <div class="text-4">
        <v-md-preview :text="post?.content" ref="preview"
                      class="dark_text_white lt-lg:important:p0"></v-md-preview>
      </div>
    </div>
    <div class="flex flex-col w-30%">
      <Profile class="mb-5"></Profile>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {getPostsById} from "~/api/post";
import type {IPostDetail} from "~/api/post";
import type {IResponse} from "~/api/http";
import {ref} from "vue";
const route = useRoute()
const id : string = String(route.params.id)
const post = ref<IPostDetail>()
const getPostDetail = async () => {
  try {
    let postRes: any = await getPostsById(id)
    let res: IResponse<IPostDetail> = postRes.data.value
    post.value = res.data
  } catch (error) {
    console.log(error);
  }
};
getPostDetail()
</script>