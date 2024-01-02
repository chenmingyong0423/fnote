<template>
  <div class="flex w-full">
    <div class="mt-10 w-5% lt-md:hidden">
      <div class="flex flex-col gap-y-3 items-center fixed ">
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white  p-2 cursor-pointer duration-200 dark:text-dtc dark_bg_gray relative"
            :class="{' hover:bg-#1e80ff': !post?.is_liked}" @click="like">
          <span
              class="i-ph:thumbs-up w-8 h-8 duration-400"
              :class="{'group-hover:scale-120 group-hover:text-white text-gray': !post?.is_liked, 'text-#1e80ff': post?.is_liked}"></span>
          <span
              class="absolute translate-x-11/10 -translate-y-2/3 bg-#1e80ff text-white text-xs rounded-full w-8 h-8 flex items-center justify-center">
              {{ post?.like_count }}
          </span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray relative"
            @click="scrollToCommentArea">
          <span
              class="i-ph-chats-duotone w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
          <span
              class="absolute transform translate-x-11/10 -translate-y-2/3 bg-#1e80ff text-white text-xs rounded-full w-8 h-8 flex items-center justify-center">
              {{ post?.comment_count }}
          </span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray relative">
          <span
              class="i-ph:share-fat-light w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
          <div
              class="share w-80% hidden absolute left-117% top--40% group-hover:block">
            <div
                class="slide-right-animation dark:text-dtc dark_bg_full_black flex flex-col gap-y-3  bg-white text-white text-xs custom_shadow_all b-rounded-4 w-150% items-center justify-center p-x-1 p-y-2 shadow-2xl shadow-black/10">
              <div class="relative w-8 h-8 group" @mouseenter="qrcodeShow = true" @mouseleave="qrcodeShow= false">
                <div class="i-bi:wechat w-full h-full text-green"></div>
                <div
                    class="qrcodeDiv absolute  left-214% top--120% w-[160px] h-[170px] p-5 bg-white b-rounded-4  duration-200 tilt-animation dark:text-dtc dark_bg_full_black"
                    :class="{'block': qrcodeShow, 'hidden ': !qrcodeShow}">
                  <div class="flex flex-col items-center align-center gap-y-3">
                    <QrcodeVue :value="`https://${domain}${path}`" :size="150" level="M"/>
                    <p class="text-black text-4 dark:text-white">微信扫一扫</p>
                  </div>
                </div>
              </div>
              <hr class="w-50% border-gray-1 b-rounded-4">
              <a class="i-bi:tencent-qq w-8 h-8 text-black dark:text-orange"
                 :href="`https://connect.qq.com/widget/shareqq/index.html?url=${domain}${path}&title=${post?.title}&pics=${post?.cover_img}`"
                 target="_blank"></a>
              <hr class="w-50% border-gray-1 b-rounded-4">
              <a class="i-bi:sina-weibo w-8 h-8 text-red"
                 :href="`https://service.weibo.com/share/share.php?sharesource=weibo&title=${post?.title}，原文链接：${domain}${path}&pic=${post?.cover_img}`"
                 target="_blank"></a>
              <hr class="w-50% border-gray-1 b-rounded-4">
              <span class="i-bi:link-45deg w-8 h-8 text-black dark:text-white" @click="copyLink"></span>
            </div>
          </div>
        </div>
        <div
            class="relative group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray">
           <span
               class="w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400 text-5 text-center">赏</span>
          <div
              class="pay slide-right-4-reword-animation dark_bg_full_black w-[640px] h-[320px] hidden absolute bg-gray-1 b-rounded-4 left-117% top--28% group-hover:block custom_shadow_all p-2">
            <div
                class="flex align-center items-center justify-center center gap-x-5">
              <img :src="code.image" width="300" height="300" :alt="code.name" v-for="(code, index) in payList"
                   :key="index">
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="w-63% ml-1% mr-1% lt-md:w-100% lt-md:mx-0%">
      <div class="bg-white mb-5 b-rounded-4 dark:text-dtc dark_bg_gray">
        <!--  文章标题  -->
        <div class="text-10 font-bold text-center p-1">{{ post?.title }}</div>
        <!--  文章 meta  -->
        <div class="flex items-center gap-x-2 text-4 justify-center p-1 text-gray-4">
          <div>阅读 {{ post?.visit_count }}</div>
        </div>
        <!--  文章内容  -->
        <div class="text-4" ref="previewRef">
          <client-only>
            <v-md-preview :text="post?.content" @copy-code-success="handleCopyCodeSuccess"
                          class="lt-lg:important:p0" :class="{'dark': isBlackMode}"
                          @change="generateAnchors"></v-md-preview>
          </client-only>
        </div>
      </div>
      <div>
        <div
            class="m-auto flex items-center justify-center align-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer duration-200 dark:text-dtc dark_bg_gray relative"
            :class="{' hover:bg-#1e80ff': !post?.is_liked}" @click="like">
          <span
              class="i-ph:thumbs-up w-8 h-8 duration-400"
              :class="{'text-gray': !post?.is_liked, 'text-#1e80ff': post?.is_liked}"></span>
          <span
              class="absolute translate-x-11/10 -translate-y-2/3 bg-#1e80ff text-white text-xs rounded-full w-8 h-8 flex items-center justify-center">
              {{ post?.like_count }}
          </span>
        </div>
      </div>
      <!-- 评论区 -->
      <div ref="comment">
        <CommentPost ref="commentPost" :comments="comments" :author="author"
                     class="mt-5 b-rounded-4 p-2 dark:text-dtc dark_bg_gray"
                     @submit="submit" @submitReply="submitReply" @submitReply2Reply="submitReply2Reply"></CommentPost>
      </div>
    </div>
    <div class="flex flex-col w-30% lt-md:hidden">
      <Profile class="mb-5"></Profile>
      <div ref="anchor">
        <Anchor :htmlContent="htmlContent" :lineIndex="lineIndex" @handleAnchorClick="handleAnchorClick"
                class="dark:text-dtc dark_bg_gray"></Anchor>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {type IPostDetail, getPostsById, likePost} from "~/api/post";
import type {IResponse, IBaseResponse, IPageData} from "~/api/http";
import {onMounted, ref} from "vue";
import {useHomeStore} from '~/store/home';

const homeStore = useHomeStore()
const isBlackMode = computed(() => homeStore.isBlackMode)

const domain = homeStore.website_info.domain;
const route = useRoute()
const path: string = route.path
const id: string = 'about-me'
const post = ref<IPostDetail>()
const author = ref<string>("")
const payList = ref<IPayInfo[]>(homeStore.pay_info || [])

const getPostDetail = async () => {
  try {
    let postRes: any = await getPostsById(id)
    let res: IResponse<IPostDetail> = postRes.data.value
    post.value = res.data
    author.value = post.value?.author || ""
  } catch (error) {
    console.log(error);
  }
};
await getPostDetail()

const handleCopyCodeSuccess = () => {
  console.log("成功")
}

const htmlContent = ref<string>("")

const generateAnchors = (text: string, html: string) => {
  htmlContent.value = html
}

const previewRef = ref<HTMLElement>()
const anchor = ref()
const anchorOriginTop = ref(0)
const lineIndex = ref("8")

const anchorScroll = () => {
  if (document.documentElement.scrollTop > anchor.value.offsetTop && anchorOriginTop.value == 0) {
    anchorOriginTop.value = anchor.value.offsetTop
    anchor.value.setAttribute('style', `position:fixed;top:88;width:${anchor.value.clientWidth}px;`);
  } else if (document.documentElement.scrollTop < anchorOriginTop.value) {
    anchorOriginTop.value = 0
    anchor.value.removeAttribute('style');
  }
}

// 用于判断是否正在滚动，滚动则不触发标题滚动监听事件
const isScrolling = ref(false)

const subscribeTitleFocus = () => {
  if (isScrolling.value) return
  // 获取当前滚动位置
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const preview = previewRef.value;
  if (!preview) return;
  const titles = preview.querySelectorAll('h1,h2,h3,h4,h5,h6')
  titles.forEach((title, _) => {
    const cur = title as HTMLElement
    if (cur.offsetTop - 60 <= scrollTop) {
      const lineIdx = cur.getAttribute("data-v-md-line");
      lineIndex.value = String(lineIdx)
      return
    }
  })
};

onMounted(() => {
  window.addEventListener('scroll', anchorScroll)
  window.addEventListener('scroll', subscribeTitleFocus)
})
onBeforeUnmount(() => {
  window.removeEventListener('scroll', anchorScroll)
  window.removeEventListener('scroll', subscribeTitleFocus)
})


const handleAnchorClick = (newLineIndex: string) => {
  isScrolling.value = true
  const preview = previewRef.value;
  if (!preview) return;
  const heading = preview.querySelector(`[data-v-md-line="${newLineIndex}"]`);
  if (heading && heading instanceof HTMLElement) {
    const top = heading.offsetTop - 60
    window.scrollTo({
      top: top,
      behavior: 'smooth',
    });

    const checkIfDone = () => {
      const currentPosition = document.documentElement.scrollTop;

      if (Math.abs(currentPosition - top) < 1 || (window.innerHeight + currentPosition) >= document.body.offsetHeight) {
        // 滚动已经完成
        lineIndex.value = newLineIndex
        isScrolling.value = false;
      } else {
        // 滚动未完成，继续检查
        requestAnimationFrame(checkIfDone);
      }
    }

    requestAnimationFrame(checkIfDone);
  } else {
    isScrolling.value = false
  }
};

// 点赞
const like = async () => {
  if (post.value?.is_liked) return
  try {
    let likeRes: any = await likePost(id)
    let res: IBaseResponse = likeRes.data.value
    if (res?.code === 200) {
      post.value!.is_liked = true
      post.value!.like_count++
    } else {
      console.log(res);
    }
  } catch (error) {
    console.log(error);
  }
}

const comment = ref()
const scrollToCommentArea = () => {
  window.scrollTo({
    top: comment.value.offsetTop - 60,
    behavior: 'smooth',
  });
};

// 二维码
const qrcodeShow = ref(false)


import {useAlertStore} from '~/store/toast';
import {
  getComments,
  type IComment,
  type ICommentReplyRequest,
  type ICommentRequest,
  submitComment, submitCommentReply
} from "~/api/comment";
import type {IPayInfo} from "~/api/config";

const toast = useAlertStore();

const copyLink = async () => {
  await navigator.clipboard.writeText(`https://${domain}${path}`);
  toast.showToast('复制成功！', 2000);
}

const comments = ref<IComment[]>([])
const initComments = async () => {
  try {
    let commentRes: any = await getComments(id)
    let res: IResponse<IPageData<IComment>> = commentRes.data.value
    if (res.code !== 200) {
      toast.showToast(res.message, 2000);
      return
    }
    comments.value = res.data?.list || []
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
}

initComments()

const submit = async (req: ICommentRequest) => {
  try {
    req.postId = id
    let commentRes: any = await submitComment(req)
    if (commentRes.data.value === null) {
      if (commentRes.error.value.statusCode == 403) {
        toast.showToast("评论模块暂未开放！", 2000);
      } else {
        toast.showToast(commentRes.error.value.statusMessage, 2000);
      }
      return
    }
    let res: IBaseResponse = commentRes.data.value
    if (res.code !== 200) {
      toast.showToast(res.message, 2000);
      return
    }
    toast.showToast("提交评论成功，待站长审核通过后将会通过邮件告知。", 3000);
    clearCommentReq()
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
}

const submitReply = async (req: ICommentReplyRequest, commentId: string) => {
  try {
    req.postId = id
    let commentRes: any = await submitCommentReply(commentId, req)
    if (commentRes.data.value === null) {
      if (commentRes.error.value.statusCode == 403) {
        toast.showToast("评论模块暂未开放！", 2000);
      } else {
        toast.showToast(commentRes.error.value.statusMessage, 2000);
      }
      return
    }
    let res: IBaseResponse = commentRes.data.value
    if (res.code !== 200) {
      toast.showToast(res.message, 2000);
      return
    }
    toast.showToast("提交评论成功，待站长审核通过后将会通过邮件告知。", 3000);
    clearCommentReplyReq()
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
}

const commentPost = ref()
const clearCommentReq = () => {
  if (commentPost.value) {
    commentPost.value.clearReq()
  }
}
const clearCommentReplyReq = () => {
  if (commentPost.value) {
    commentPost.value.clearReplyReq()
  }
}

const submitReply2Reply = async (req: ICommentReplyRequest, commentId: string) => {
  try {
    req.postId = id
    let commentRes: any = await submitCommentReply(commentId, req)
    if (commentRes.data.value === null) {
      if (commentRes.error.value.statusCode == 403) {
        toast.showToast("评论模块暂未开放！", 2000);
      } else {
        toast.showToast(commentRes.error.value.statusMessage, 2000);
      }
      return
    }
    let res: IBaseResponse = commentRes.data.value
    if (res.code !== 200) {
      toast.showToast(res.message, 2000);
      return
    }
    toast.showToast("提交评论成功，待站长审核通过后将会通过邮件告知。", 3000);
    clearReply2ReplyReq()
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
}

let description = post.value?.summary
if (post.value?.meta_description) {
  description = post.value?.meta_description
}
const clearReply2ReplyReq = () => {
  if (commentPost.value) {
    commentPost.value.clearReply2ReplyReq()
  }
}

useHead({
  title: `${post.value?.title} - ${homeStore.seo_meta_config.title}`,
  meta: [
    {name: 'description', content: description},
    {name: 'keywords', content: homeStore.seo_meta_config.keywords},
    {name: 'author', 'content': homeStore.seo_meta_config.author},
    {name: 'robots', 'content': homeStore.seo_meta_config.robots},
  ],
  link: [
    {rel: 'icon', type: 'image/x-icon', href: homeStore.website_info.icon},
  ]
})
useSeoMeta({
  ogTitle: `${post.value?.title} - ${homeStore.seo_meta_config.title}`,
  ogDescription: description,
  ogImage: '',
  twitterCard: 'summary'
})
</script>

<style scoped>
.dark :deep(.v-md-pre-wrapper) {
  background-color: rgba(10, 0, 0, 0.1) !important;
}

.dark :deep(.v-md-pre-wrapper) {
  background-color: rgba(10, 0, 0, 0.1) !important;
}

.dark :deep(code) {
  color: white !important;
}

.dark :deep(.line-numbers-mode:after) {
  background-color: rgba(10, 0, 0, 0.1) !important;
  border: 0 !important;
}

/* 根据需要定制不同代码语言或元素的样式 */
.dark :deep(.hljs-keyword, .hljs-selector-tag, .hljs-literal) {
  color: #ff7b72 !important;
}

.dark :deep(.hljs-string) {
  color: #a5d6ff !important;
}

.dark :deep(.hljs-title) {
  color: #a5d6ff !important;
}

.dark :deep(.hljs-type) {
  color: #cc880a !important;
}

.dark :deep(.github-markdown-body table tr) {
  background-color: rgba(10, 0, 0, 0.1) !important;

}

.dark :deep(.github-markdown-body blockquote) {
  border-color: #334a61 !important;
}


.share::before {
  content: '';
  position: absolute;
  left: -20px;
  top: 25%;
  border-width: 10px;
  border-style: solid;
  border-color: transparent #b7bbc4 transparent transparent;
  transform: translateY(-50%);
}

.qrcodeDiv::before {
  content: '';
  position: absolute;
  left: -20px;
  top: 25%;
  border-width: 10px;
  border-style: solid;
  border-color: transparent #b7bbc4 transparent transparent;
  transform: translateY(-50%);
}


@keyframes tilt {
  0%, 100% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(-10deg);
  }
  50% {
    transform: rotate(10deg);
  }
  75% {
    transform: rotate(-5deg);
  }
}

.tilt-animation {
  animation: tilt 0.5s ease-in-out;
}


@keyframes slideRight {
  from {
    transform: translateX(-30%);
  }
  to {
    transform: translateX(0);
  }
}

.slide-right-animation {
  animation: slideRight 0.3s ease-out;
}

@keyframes slideRight4Reword {
  from {
    transform: translateX(-10%);
  }
  to {
    transform: translateX(0);
  }
}

.slide-right-4-reword-animation {
  animation: slideRight4Reword 0.3s ease-out;
}

.pay:before {
  content: '';
  position: absolute;
  left: -20px;
  top: 10%;
  border-width: 10px;
  border-style: solid;
  border-color: transparent #b7bbc4 transparent transparent;
  transform: translateY(-50%);
}
</style>
