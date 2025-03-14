<template>
  <div class="flex mb-5" v-for="(rpy, index) in replies" :key="index">
    <div class="w-8% min-h-[100px] flex justify-center lt-md:w-15%">
      <img
        :src="rpy.picture"
        alt=""
        class="w-12 h-12 border-rounded-50% cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
        v-if="rpy.picture != ''"
      />
      <div
        class="i-ph-user-circle-duotone w-12 h-12 border-rounded-50% lt-lg:mr0 text-gray-4"
        v-else
      ></div>
    </div>
    <div class="w-91% ml-1% flex flex-col lt-md:w-84%">
      <div
        class="text-gray-4 h-[55px] line-height-[35px] flex gap-x-2 lt-md:line-height-[25px]"
      >
        <div class="flex gap-x-2 lt-md:flex-col">
          <div class="flex gap-x-2 lt-md:gap-x-1">
            <a
              v-if="rpy.website !== ''"
              :href="rpy.website"
              target="_blank"
              class="text-#1E80FF lt-md:truncate"
              >{{
                rpy.name === props.author ? `${rpy.name}[作者]` : rpy.name
              }}</a
            >
            <span v-else class="text-#1E80FF lt-md:truncate">{{
              rpy.name === props.author ? `${rpy.name}[作者]` : rpy.name
            }}</span>
            <span>回复</span>
            <span class="text-#1E80FF lt-md:truncate">{{
              rpy.reply_to === props.author
                ? `${rpy.reply_to}[作者]`
                : rpy.reply_to
            }}</span>
          </div>
          <span>
            发表于
            {{
              $dayjs(rpy.reply_time * 1000).format("YYYY-MM-DD HH:mm:ss")
            }}</span
          >
        </div>
        <Button
          name="回复"
          class="w-15 h-8 line-height-8 hover:bg-gray-1 ml-auto"
          @click="
            activeCommentIndex = rpy.id;
            console.log(rpy.id);
          "
        ></Button>
      </div>
      <div>
        <div :data-theme="isBlackMode ? 'dark' : 'light'">
          <MDC
              :value="rpy.content"
              class="markdown-body lt-lg:important:p0"
          />
        </div>
      </div>
      <div
        class="custom_bg_gray h-10 line-height-10 pl-4 truncate dark:text-dtc dark_bg_gray"
        v-if="rpy.replied_content !== '' && rpy.replied_content !== undefined"
      >
        {{ rpy.replied_content }}
      </div>
      <div v-if="activeCommentIndex === rpy.id">
        <CommentForm
          class="m-auto"
          @submit="submitReply"
          ref="replyReplyForm"
          :commentId="rpy.id"
        ></CommentForm>
        <Button
          name="取消"
          class="w-15 h-8 line-height-8 hover:bg-gray-1 m-auto"
          @click="activeCommentIndex = ''"
        ></Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useHomeStore } from "~/store/home";
import type {
  ICommentReplyRequest,
  ICommentRequest,
  IReply,
} from "~/api/comment";
import type { PropType } from "vue";

const props = defineProps({
  replies: {
    type: Array as PropType<IReply[]>,
    default: () => [],
  },
  author: {
    type: String,
    default: "",
  },
  commentId: {
    type: String,
    required: true,
  },
});

const replies = computed(() => {
  return props.replies.map((rpy) => {
    const newRpy = { ...rpy }; // 创建一个新对象
    if (newRpy.reply_to_id !== "") {
      const replied = props.replies.find((r) => r.id === newRpy.reply_to_id);
      if (replied) {
        newRpy.replied_content = replied.content;
      }
    }
    return newRpy;
  });
});

const homeStore = useHomeStore();
const isBlackMode = computed(() => homeStore.isBlackMode);
const emit = defineEmits(["submitReply2Reply"]);

const activeCommentIndex = ref("");
const submitReply = (req: ICommentRequest, replyToId: string) => {
  const req4Reply: ICommentReplyRequest = {
    ...req,
    replyToId: replyToId,
  };
  emit("submitReply2Reply", req4Reply, props.commentId);
};

const replyReplyForm = ref();
const clearReplyReq = () => {
  if (replyReplyForm.value) {
    replyReplyForm.value[0].clearReq();
    activeCommentIndex.value = "";
  }
};

defineExpose({
  clearReplyReq,
});
</script>

<style scoped>
.markdown-body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 45px;
}

@media (max-width: 767px) {
  .markdown-body {
    padding: 15px;
  }
}

.markdown-body :deep(a) {
  color: black !important;
}

.dark .markdown-body :deep(a) {
  color: #ffffffb2 !important;
}
</style>
