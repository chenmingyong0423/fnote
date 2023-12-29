<template>
  <div class="flex mb-5" v-for="(rpy, index) in replies" :key="index">
    <div class="w-8% min-h-[180px] ml-1% flex justify-center">
      <img :src="generateAvatar(rpy.email)" alt=""
           class="w-12 h-12 border-rounded-50%  cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
           v-if="rpy.email != ''">
      <div class="i-ph-user-circle-duotone w-12 h-12 border-rounded-50%  lt-lg:mr0 text-gray-4" v-else></div>
    </div>
    <div class="w-91% flex flex-col ">
      <div class="text-gray-4 h-[55px] line-height-[35px] flex">
        <a v-if="rpy.website !== ''" :href="rpy.website" target="_blank" class="text-#1E80FF">{{
            rpy.name === props.author ? `${rpy.name}[作者]` : rpy.name
          }}</a>
        <span v-else class="text-#1E80FF">{{
            rpy.name === props.author ? `${rpy.name}[作者]` : rpy.name
          }}</span>
        <span> 发表于 {{ $dayjs(rpy.reply_time * 1000).format('YYYY-MM-DD HH:mm:ss') }}</span>
        <Button name="回复" class="hover:bg-gray-1 ml-auto" @click="activeCommentIndex = rpy.id"></Button>
      </div>
      <div>
        <v-md-preview :text="rpy.content"
                      class="lt-lg:important:p0" :class="{'dark': isBlackMode}"
        ></v-md-preview>
      </div>
      <div class="custom_bg_gray h-10 line-height-10 pl-4 truncate dark:text-dtc dark_bg_gray" v-if="rpy.replied_content !== '' && rpy.replied_content !== undefined">
        {{rpy.replied_content}}
      </div>
      <div v-if="activeCommentIndex === rpy.id">
        <CommentForm class="m-auto" @submit="submitReply" ref="commentReplyForm" :commentId="rpy.id"
        ></CommentForm>
        <Button name="取消" class="hover:bg-gray-1 m-auto" @click="activeCommentIndex = ''"></Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import CryptoJS from "crypto-js";
import {useHomeStore} from "~/store/home";
import type {ICommentReplyRequest, ICommentRequest, IReply} from "~/api/comment";
import type {PropType} from "vue";

const props = defineProps({
  replies: {
    type: Array as PropType<IReply[]>,
    default: () => []
  },
  author: {
    type: String,
    default: ''
  },
  commentId: {
    type: String,
    required: true
  }
})

const replies = computed(() => {
  return props.replies.map((rpy) => {
    const newRpy = { ...rpy }; // 创建一个新对象
    if (newRpy.reply_to_id !== '') {
      const replied = props.replies.find((r) => r.id === newRpy.reply_to_id);
      if (replied) {
        newRpy.replied_content = replied.content;
      }
    }
    return newRpy;
  });
})


const info = useHomeStore()
const isBlackMode = computed(() => info.isBlackMode)
const generateAvatar = (email: string) => {
  return "https://1.gravatar.com/avatar/" + CryptoJS.MD5(email.trim().toLowerCase()).toString()
}
const emit = defineEmits(['submitReply2Reply']);

const activeCommentIndex = ref('')
const submitReply = (req: ICommentRequest, replyToId: string) => {
  const req4Reply : ICommentReplyRequest= {
    ...req,
    replyToId: replyToId,
  }
  emit('submitReply2Reply', req4Reply, props.commentId)
}

const commentReplyForm = ref()
const clearReplyReq = () => {
  if (commentReplyForm.value) {
    commentReplyForm.value[0].clearReq()
  }
}

defineExpose({
  clearReplyReq
});
</script>