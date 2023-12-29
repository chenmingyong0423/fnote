<template>
  <div class="bg-white p-2">
    <div class="h-[50px] line-height-[50px] pl-3 custom_bottom_border_gray mb-2 flex items-center gap-x-1">
      <span class="i-ph-chats-duotone text-2xl block text-black dark:text-dtc"></span> <span class="text-2xl">评论</span>
    </div>
    <div class="mb-2 custom_bottom_border_gray p-1">
      <CommentForm class="m-auto" @submit="submit" ref="commentForm"></CommentForm>
    </div>
    <div class="flex mb-5" v-for="(comment, index) in props.comments" :key="index">
      <div class="w-8% min-h-[180px] ml-1% flex justify-center">
        <img :src="generateAvatar(comment.email)" alt=""
             class="w-12 h-12 border-rounded-50%  cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
             v-if="comment.email != ''">
        <div class="i-ph-user-circle-duotone w-12 h-12 border-rounded-50%  lt-lg:mr0 text-gray-4" v-else></div>
      </div>
      <div class="w-91% flex flex-col ">
        <div class="text-gray-4 h-[55px] line-height-[35px] flex">
          <a v-if="comment.website !== ''" :href="comment.website" target="_blank" class="text-#1E80FF">{{
              comment.username === props.author ? `${comment.username}[作者]` : comment.username
            }}</a>
          <span v-else class="text-#1E80FF">{{
              comment.username === props.author ? `${comment.username}[作者]` : comment.username
            }}</span>
          <span> 发表于 {{ $dayjs(comment.comment_time * 1000).format('YYYY-MM-DD HH:mm:ss') }}</span>
          <Button name="回复" class="hover:bg-gray-1 ml-auto" @click="activeCommentIndex = comment.id"></Button>
        </div>
        <div>
          <v-md-preview :text="comment.content"
                        class="lt-lg:important:p0" :class="{'dark': isBlackMode}"
          ></v-md-preview>
        </div>
        <div>
          <CommentReply :replies="comment.replies" :author="props.author" :commentId="comment.id" @submitReply2Reply="submitReply2Reply" ref="replyList"></CommentReply>
        </div>
        <div v-if="activeCommentIndex === comment.id">
          <CommentForm class="m-auto" @submit="submitReply" ref="commentReplyForm" :commentId="comment.id"
                       ></CommentForm>
          <Button name="取消" class="hover:bg-gray-1 m-auto" @click="activeCommentIndex = ''"></Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useHomeStore} from "~/store/home";
import type {IComment, ICommentReplyRequest, ICommentRequest} from "~/api/comment";
import type {PropType} from "vue";
import CryptoJS from 'crypto-js'

const info = useHomeStore()
const isBlackMode = computed(() => info.isBlackMode)

const props = defineProps({
  comments: {
    type: Array as PropType<IComment[]>,
    default: () => []
  },
  author: {
    type: String,
    default: ''
  }
})

const commentForm = ref()

const clearReq = () => {
  console.log(commentForm.value)
  if (commentForm.value) {
    commentForm.value.clearReq()
  }
}

const commentReplyForm = ref()
const clearReplyReq = () => {
  if (commentReplyForm.value) {
    commentReplyForm.value[0].clearReq()
  }
}



const generateAvatar = (email: string) => {
  return "https://1.gravatar.com/avatar/" + CryptoJS.MD5(email.trim().toLowerCase()).toString()
}

const emit = defineEmits(['submit', 'submitReply', 'submitReply2Reply']);

const submit = (req: ICommentRequest, commentId: string) => {
  emit('submit', req)
}


const activeCommentIndex = ref('')
const submitReply = (req: ICommentRequest, commentId: string) => {
  const req4Reply : ICommentReplyRequest= {
    ...req,
  }
  emit('submitReply', req4Reply, commentId)
}

const submitReply2Reply = (req: ICommentReplyRequest, replyToId: string) => {
  emit('submitReply2Reply', req, replyToId)
}

const replyList = ref()
const clearReply2ReplyReq = () => {
  if (replyList.value) {
    replyList.value[0].clearReplyReq()
  }
}

defineExpose({
  clearReq,
  clearReplyReq,
  clearReply2ReplyReq
});
</script>

<style scoped>
/deep/ .github-markdown-body {
  padding: 0;
}
</style>