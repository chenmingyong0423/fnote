<template>
  <div>
    <div>
      <textarea v-if="!isPreview" rows="10"
                class="w-full custom_border_gray bg-#F9F9F9 outline-none focus:custom_border_1E80FF b-rounded-2 p-2 box-border mb-3"
                v-model="commentReq.content" maxlength="200"></textarea>
      <div v-else>
        <div class="font-bold flex items-center gap-x-2">
          <span class="i-ph:eye w-8 h-8 block text-black"></span> <span>预览中...</span>
        </div>
        <br>
        <v-md-preview :text="commentReq.content"
                      class="lt-lg:important:p0 max-h-[170px] h-[170px] overflow-x-auto text-nowrap custom_border_gray p-2 b-rounded-2"
                      :class="{'dark': isBlackMode}"
        ></v-md-preview>
      </div>
    </div>
    <div>
      <div class="h-10 text-center line-height-10">
        个人信息
      </div>
      <div class="flex justify-between items-center w-full mb-3">
        <div class="w-5%">
          <img v-if="pic != ''" :src="pic" alt=""
               class="w-auto h-12 border-rounded-50%  cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0">
          <div class="i-ph-user-circle-duotone w-auto h-12 border-rounded-50%  lt-lg:mr0 text-gray-4" v-else></div>
        </div>
        <div class="w-30% relative">
          <input type="text" placeholder="* 昵称" v-model="commentReq.username"
                 class="w-full outline-none custom_border_gray bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border"
                 @focusin="showUsernameTip=true"
                 @focusout="showUsernameTip=false">
          <span
              class="popup-text absolute bg-#555 left-25% bottom-125% text-white p-2 b-rounded-2 h-[20px] animated-fadeIn"
              v-if="showUsernameTip">您的昵称？</span>
        </div>
        <div class="w-30% relative">
          <input type="text" placeholder="* 邮箱" v-model="commentReq.email"
                 class="w-full outline-none custom_border_gray bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border"
                 @focusin="showEmailTip=true"
                 @focusout="calculateMD54Email">
          <span
              class="popup-text absolute bg-#555 left-25% bottom-125% text-white p-2 b-rounded-2 h-[20px] animated-fadeIn"
              v-if="showEmailTip">用于接收通知。</span>
        </div>
        <div class="w-30% relative">
          <input type="text" placeholder="个人站点" v-model="commentReq.website"
                 class="w-full outline-none custom_border_gray bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border"
                 @focusin="showWebsiteTip=true"
                 @focusout="showWebsiteTip=false">
          <span
              class="popup-text absolute bg-#555 left-25% bottom-125% text-white p-2 b-rounded-2 h-[20px] animated-fadeIn"
              v-if="showWebsiteTip">不许打广告哟！</span>
        </div>
      </div>
      <div class="flex justify-end gap-x-2">
        <Button :name="isPreview ? '编辑' : '预览'" class="bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200"
                @click="isPreview = !isPreview"></Button>
        <Button name="提交" class="bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200" @click="submit"></Button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import {useHomeStore} from "~/store/home";
import type {ICommentRequest} from "~/api/comment";
import CryptoJS from 'crypto-js'

const showUsernameTip = ref(false)
const showEmailTip = ref(false)
const showWebsiteTip = ref(false)
const pic = ref<string>('')


const commentReq = ref<ICommentRequest>({
  postId: "",
  username: "",
  email: "",
  website: "",
  content: "",
})
const info = useHomeStore()
const isBlackMode = computed(() => info.isBlackMode)
const isPreview = ref(false)

const emit = defineEmits(['submit']);
const submit = () => {
  if (commentReq.value.content === "") {
    alert("评论内容不能为空")
    return
  } else if (commentReq.value.username === "") {
    alert("昵称不能为空")
    return
  } else if (commentReq.value.email === "") {
    alert("邮箱不能为空")
    return
  }
  const deepCopyReq: ICommentRequest = JSON.parse(JSON.stringify(commentReq.value));
  emit("submit", deepCopyReq)
}

const calculateMD54Email = () => {
  showEmailTip.value = false
  console.log(commentReq.value.email)
  if (commentReq.value.email !== "") {
    console.log(123)
    pic.value = "https://1.gravatar.com/avatar/" + CryptoJS.MD5(commentReq.value.email.trim().toLowerCase()).toString()
  } else {
    pic.value = ''
  }
}

</script>

<style scoped>
.popup-text:before {
  content: '';
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
