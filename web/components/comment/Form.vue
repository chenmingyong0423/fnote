<template>
  <div class="p-y-2">
    <div>
      <textarea
        v-if="!isPreview"
        rows="10"
        class="w-full custom_border_gray bg-#F9F9F9 outline-none focus:custom_border_1E80FF b-rounded-2 p-2 box-border mb-3 dark:text-dtc dark_bg_gray"
        v-model="commentReq.content"
        maxlength="200"
      ></textarea>
      <div v-else>
        <div class="font-bold flex items-center gap-x-2">
          <span class="i-ph:eye w-8 h-8 block text-black"></span>
          <span>预览中...</span>
        </div>
        <br />
        <div :data-theme="isBlackMode ? 'dark' : 'light'">
          <MDC
            :value="commentReq.content"
            class="markdown-body lt-lg:important:p0 max-h-[170px] h-[170px] overflow-x-auto text-nowrap custom_border_gray p-2 b-rounded-2"
          />
        </div>
      </div>
    </div>
    <div>
      <div class="h-10 text-center line-height-10">个人信息</div>
      <div
        class="flex justify-between items-center w-full mb-3 lt-md:flex-col lt-md:gap-y-2 lt-md:justify-center"
      >
        <div class="w-5%">
          <img
            v-if="pic != ''"
            :src="pic"
            alt=""
            class="w-auto h-12 border-rounded-50% cursor-pointer hover:rotate-360 ease-out duration-1000"
          />
          <div
            class="i-ph-user-circle-duotone w-auto h-12 border-rounded-50% text-gray-4"
            v-else
          ></div>
        </div>
        <div class="w-30% relative lt-md:w-100%">
          <input
            type="text"
            placeholder="* 昵称"
            v-model="commentReq.username"
            class="w-full outline-none custom_border_gray bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
            @focusin="showUsernameTip = true"
            @focusout="showUsernameTip = false"
          />
          <span
            class="popup-text absolute bg-#555 left-25% bottom-125% text-white p-2 b-rounded-2 h-[20px] animated-fadeIn"
            v-if="showUsernameTip"
            >您的昵称？</span
          >
        </div>
        <div class="w-30% relative lt-md:w-100%">
          <input
            type="text"
            placeholder="* 邮箱"
            v-model="commentReq.email"
            class="w-full outline-none custom_border_gray bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
            @focusin="showEmailTip = true"
            @focusout="calculateMD54Email"
          />
          <span
            class="popup-text absolute bg-#555 left-25% bottom-125% text-white p-2 b-rounded-2 h-[20px] animated-fadeIn"
            v-if="showEmailTip"
            >用于接收通知。</span
          >
        </div>
        <div class="w-30% relative lt-md:w-100%">
          <input
            type="text"
            placeholder="个人站点，以 https:// 开头"
            v-model="commentReq.website"
            class="w-full outline-none custom_border_gray bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
            @focusin="showWebsiteTip = true"
            @focusout="showWebsiteTip = false"
          />
          <span
            class="popup-text absolute bg-#555 left-25% bottom-125% text-white p-2 b-rounded-2 h-[20px] animated-fadeIn"
            v-if="showWebsiteTip"
            >不许打广告哟！</span
          >
        </div>
      </div>
      <div class="flex gap-x-2 justify-center">
        <Button
          name="清空"
          class="w-15 h-8 line-height-8 bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200"
          @click="clearReq"
        ></Button>
        <Button
          :name="isPreview ? '编辑' : '预览'"
          class="w-15 h-8 line-height-8 bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200"
          @click="isPreview = !isPreview"
        ></Button>
        <Button
          name="提交"
          class="w-15 h-8 line-height-8 bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200"
          @click="submit"
        ></Button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useHomeStore } from "~/store/home";
import type { ICommentRequest } from "~/api/comment";
import CryptoJS from "crypto-js";
import { useAlertStore } from "~/store/toast";

const props = defineProps({
  commentId: {
    type: String,
    default: "",
    required: false,
  },
});

const showUsernameTip = ref(false);
const showEmailTip = ref(false);
const showWebsiteTip = ref(false);
const pic = ref<string>("");
const commentId = ref<string>(props.commentId);

const commentReq = ref<ICommentRequest>({
  postId: "",
  username: "",
  email: "",
  website: "",
  content: "",
});
const info = useHomeStore();
const isBlackMode = computed(() => info.isBlackMode);
const isPreview = ref(false);

const toast = useAlertStore();

const emit = defineEmits(["submit"]);

import { isValidEmail } from "~/utils/email";

const submit = () => {
  if (commentReq.value.content === "") {
    toast.showToast("评论内容不能为空！", 1000);
    return;
  } else if (commentReq.value.username === "") {
    toast.showToast("昵称不能为空！", 1000);
    return;
  } else if (commentReq.value.email === "") {
    toast.showToast("邮箱不能为空！", 1000);
    return;
  }
  if (
    commentReq.value.website !== "" &&
    !commentReq.value.website?.startsWith(`https://`)
  ) {
    toast.showToast("个人站点格式不正确！", 1000);
    return;
  }
  if (!isValidEmail(commentReq.value.email)) {
    toast.showToast("邮箱格式不正确！", 1000);
    return;
  }

  const deepCopyReq: ICommentRequest = JSON.parse(
    JSON.stringify(commentReq.value),
  );
  emit("submit", deepCopyReq, commentId.value);
};

const clearReq = () => {
  commentReq.value = {
    postId: "",
    username: "",
    email: "",
    website: "",
    content: "",
  };
  isPreview.value = false;
};

defineExpose({
  clearReq,
});

const calculateMD54Email = () => {
  showEmailTip.value = false;
  if (commentReq.value.email !== "") {
    pic.value =
      "https://1.gravatar.com/avatar/" +
      CryptoJS.MD5(commentReq.value.email.trim().toLowerCase()).toString();
  } else {
    pic.value = "";
  }
};
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

.popup-text:before {
  content: "";
  position: absolute;
  left: 50%;
  bottom: -20px;
  border-width: 10px;
  border-style: solid;
  border-color: #555 transparent transparent transparent;
  transform: translateX(-50%);
}

/* 定义关键帧动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* 初始状态 - 设置标签透明 */
.animated-fadeIn {
  opacity: 0;
  animation: fadeIn 0.5s ease-in-out forwards; /* 应用关键帧动画 */
}
</style>
