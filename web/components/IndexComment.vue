<template>
  <div class="flex flex-col  bg-white bg-#fff p-5 b-rounded-4 dark_bg_gray dark:text-dtc">
    <div class="text-6 border-b-1 border-b-gray-2 border-b-solid p-x-1 p-y-2">最新评论</div>
    <div class="border-b-1 border-b-gray-2 border-b-solid p-b-2 p-2" v-for="(item, index) in comments" :key="index">
      <div class="flex h-15 my-2">
        <div>
          <img :src="item.picture" alt=""
               class="w-15 h-15 border-rounded-50%  cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
               v-if="item.picture">
          <div class="i-ph-user-circle-duotone w-15 h-15 border-rounded-50%  lt-lg:mr0 text-gray-4" v-else></div>
        </div>
        <div class="flex flex-col justify-center items-start ml-3">
          <span class="text-5">{{ item.name }}</span>
          <span class="text-gray-5">{{ $dayjs(item.create_time).format('YYYY-MM-DD HH:mm:ss') }}</span>
        </div>
      </div>
      <div class="">
        <div class="p-y-2 truncate">
          {{ item.content }}
        </div>
        <div class="flex gap-2 items-center text-gray-5 p-y-1 cursor-pointer hover:bg-green-1 "
             @click="router.push('/posts/' + item.post_id)">
          <span class="i-ph-notebook-light "></span>
          <span>{{ item.post_title }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import {getLatestComments} from "~/server/api/comment"
import type {ILatestComment,} from "~/server/api/comment"
import type {IResponse, IListData} from "~/server/api/http";

const router = useRouter()
const comments = ref([] as ILatestComment[]);
const commentsInfos = async () => {
  try {
    let postRes: any = await getLatestComments()
    let res: IResponse<IListData<ILatestComment>> = postRes.data.value
    comments.value = res.data?.list || []
  } catch (error) {
    console.log(error);
  }
};
commentsInfos()
</script>