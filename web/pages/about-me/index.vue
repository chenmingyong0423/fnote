<template>
  <div class="flex w-full slide-up">
    <div class="mt-10 w-5% lt-md:hidden" v-if="existPost">
      <div class="flex flex-col gap-y-3 items-center fixed">
        <div
          class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer duration-200 dark:text-dtc dark_bg_gray relative"
          :class="{ ' hover:bg-#1e80ff': !post?.is_liked }"
          @click="like"
        >
          <span
            class="i-ph:thumbs-up w-8 h-8 duration-400"
            :class="{
              'group-hover:scale-120 group-hover:text-white text-gray':
                !post?.is_liked,
              'text-#1e80ff': post?.is_liked,
            }"
          ></span>
          <span
            class="absolute translate-x-11/10 -translate-y-2/3 bg-#1e80ff text-white text-xs rounded-full w-8 h-8 flex items-center justify-center"
          >
            {{ post?.like_count }}
          </span>
        </div>
        <div
          class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray relative"
          @click="scrollToCommentArea"
        >
          <span
            class="i-ph-chats-duotone w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"
          ></span>
          <span
            class="absolute transform translate-x-11/10 -translate-y-2/3 bg-#1e80ff text-white text-xs rounded-full w-8 h-8 flex items-center justify-center"
          >
            {{ post?.comment_count }}
          </span>
        </div>
        <div
          class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray relative"
        >
          <span
            class="i-ph:share-fat-light w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"
          ></span>
          <div
            class="share w-80% hidden absolute left-117% top--40% group-hover:block"
          >
            <div
              class="slide-right-animation dark:text-dtc dark_bg_full_black flex flex-col gap-y-3 bg-white text-white text-xs custom_shadow_all b-rounded-4 w-150% items-center justify-center p-x-1 p-y-2 shadow-2xl shadow-black/10"
            >
              <div
                class="relative w-8 h-8 group"
                @mouseenter="qrcodeShow = true"
                @mouseleave="qrcodeShow = false"
              >
                <div class="i-bi:wechat w-full h-full text-green"></div>
                <div
                  class="qrcodeDiv absolute left-214% top--120% w-[160px] h-[170px] p-5 bg-white b-rounded-4 duration-200 tilt-animation dark:text-dtc dark_bg_full_black"
                  :class="{ block: qrcodeShow, 'hidden ': !qrcodeShow }"
                >
                  <div class="flex flex-col items-center align-center gap-y-3">
                    <QrcodeVue :value="link" :size="150" level="M" />
                    <p class="text-black text-4 dark:text-white">微信扫一扫</p>
                  </div>
                </div>
              </div>
              <hr class="w-50% border-gray-1 b-rounded-4" />
              <a
                class="i-bi:tencent-qq w-8 h-8 text-black dark:text-orange"
                :href="`https://connect.qq.com/widget/shareqq/index.html?url=${link}&title=${post?.title}&pics=${post?.cover_img}`"
                target="_blank"
              ></a>
              <hr class="w-50% border-gray-1 b-rounded-4" />
              <a
                class="i-bi:sina-weibo w-8 h-8 text-red"
                :href="`https://service.weibo.com/share/share.php?sharesource=weibo&title=${post?.title}，原文链接：${link}&pic=${post?.cover_img}`"
                target="_blank"
              ></a>
              <hr class="w-50% border-gray-1 b-rounded-4" />
              <span
                class="i-bi:link-45deg w-8 h-8 text-black dark:text-white"
                @click="copyLink"
              ></span>
            </div>
          </div>
        </div>
        <div
          class="relative group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray"
        >
          <span
            class="w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400 text-5 text-center"
            >赏</span
          >
          <div
            class="pay slide-right-4-reword-animation dark_bg_full_black w-[640px] h-[320px] hidden absolute bg-gray-1 b-rounded-4 left-117% top--28% group-hover:block custom_shadow_all p-2"
          >
            <div
              class="flex align-center items-center justify-center center gap-x-5"
            >
              <img
                :src="code.image"
                width="300"
                height="300"
                :alt="code.name"
                v-for="(code, index) in payList"
                :key="index"
              />
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
        <div
          class="flex items-center gap-x-2 text-4 justify-center p-1 text-gray-4"
        >
          <div>阅读 {{ post?.visit_count }}</div>
        </div>
        <!--  文章内容  -->
        <div
          class="text-4"
          ref="previewRef"
          :data-theme="isBlackMode ? 'dark' : 'light'"
        >
          <MDCRenderer
            :body="mdData.body"
            :data="mdData.data"
            class="markdown-body lt-lg:important:p0"
            tag="article"
          />
        </div>
      </div>
      <div v-if="existPost">
        <div
          class="m-auto flex items-center justify-center align-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer duration-200 dark:text-dtc dark_bg_gray relative"
          :class="{ ' hover:bg-#1e80ff': !post?.is_liked }"
          @click="like"
        >
          <span
            class="i-ph:thumbs-up w-8 h-8 duration-400"
            :class="{
              'text-gray': !post?.is_liked,
              'text-#1e80ff': post?.is_liked,
            }"
          ></span>
          <span
            class="absolute translate-x-11/10 -translate-y-2/3 bg-#1e80ff text-white text-xs rounded-full w-8 h-8 flex items-center justify-center"
          >
            {{ post?.like_count }}
          </span>
        </div>
      </div>
      <!-- 评论区 -->
      <div ref="comment" v-if="existPost">
        <CommentPost
          ref="commentPost"
          :comments="comments"
          :author="author"
          class="mt-5 b-rounded-4 p-2 dark:text-dtc dark_bg_gray"
          @submit="submit"
          @submitReply="submitReply"
          @submitReply2Reply="submitReply2Reply"
        ></CommentPost>
      </div>
    </div>
    <div class="flex flex-col w-30% lt-md:hidden">
      <Profile class="mb-5"></Profile>
      <div ref="anchor" v-if="existPost">
        <Anchor
          :toc="mdData.toc.links"
          :lineIndex="lineIndex"
          @handleAnchorClick="handleAnchorClick"
          class="dark:text-dtc dark_bg_gray"
        ></Anchor>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { type IPostDetail, getPostsById, likePost } from "~/api/post";
import type { IResponse, IBaseResponse, IPageData } from "~/api/http";
import { onMounted, ref } from "vue";
import { useHomeStore } from "~/store/home";

const homeStore = useHomeStore();
const isBlackMode = computed(() => homeStore.isBlackMode);
const configStore = useConfigStore();

const route = useRoute();
const path: string = route.path;
const id: string = "about-me";
const post = ref<IPostDetail>({
  sug: "",
  author: "",
  title: "关于我",
  summary: "",
  cover_img: "",
  category: [],
  tags: [],
  like_count: 0,
  comment_count: 0,
  visit_count: 0,
  sticky_weight: 0,
  created_at: 0,
  content: "暂无介绍。",
  meta_description: "",
  meta_keywords: "",
  word_count: 0,
  updated_at: 0,
  is_liked: false,
});
const author = ref<string>("");
const payList = ref<IPayInfo[]>(configStore.pay_info || []);
const existPost = ref<boolean>(false);
const runtimeConfig = useRuntimeConfig();

const apiHost = runtimeConfig.public.apiHost;

// 让 Nuxt 在服务器端请求数据，并在客户端复用
const { data: post2 } = await useAsyncData(`about-me`, async () => {
  let postRes: any = await getPostsById(id);
  let res: IResponse<IPostDetail> = postRes.data.value;

  if (res && res.data) {
    res.data.content = res.data.content.replace(
      /\]\(\/static\//g,
      `](${apiHost}/static/`,
    );
    return res.data;
  }
  return null;
});

if (post2.value !== null) {
  existPost.value = true;
  post.value = post2.value;
}

// 解析 Markdown 数据
const { data: mdData } = await useAsyncData(`markdown-${id}`, () =>
  post.value
    ? parseMarkdown(post.value.content, {
        toc: {
          depth: 5,
        },
      })
    : null,
);

const link = ref("");

const previewRef = ref<HTMLElement>();
const anchor = ref();
const anchorOriginTop = ref(0);
const lineIndex = ref("8");

const anchorScroll = () => {
  if (
    document.documentElement.scrollTop > anchor.value.offsetTop &&
    anchorOriginTop.value == 0
  ) {
    anchorOriginTop.value = anchor.value.offsetTop;
    anchor.value.setAttribute(
      "style",
      `position:fixed;top:88;width:${anchor.value.clientWidth}px;`,
    );
  } else if (document.documentElement.scrollTop < anchorOriginTop.value) {
    anchorOriginTop.value = 0;
    anchor.value.removeAttribute("style");
  }
};

let titleOffsets: { id: string; top: number }[] = [];

onMounted(() => {
  const preview = previewRef.value;
  if (!preview) return;

  // 缓存所有标题元素
  const titles: HTMLElement[] = Array.from(
    preview.querySelectorAll("h1,h2,h3,h4,h5,h6"),
  );

  // 预存 offsetTop，避免 scroll 事件中重复计算
  titleOffsets = titles.map((el) => ({
    id: el.getAttribute("id") || "",
    top: el.offsetTop - 20,
  }));

  window.addEventListener("scroll", anchorScroll);
  window.addEventListener("scroll", subscribeTitleFocus);
  link.value = window.location.href;
});

const subscribeTitleFocus = () => {
  const scrollTop =
    document.documentElement.scrollTop || document.body.scrollTop;

  if (titleOffsets.length === 0) return;

  // 使用二分查找定位当前标题
  let left = 0,
    right = titleOffsets.length - 1;
  let targetIdx = 0;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);
    if (titleOffsets[mid].top <= scrollTop) {
      targetIdx = mid;
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  const newHash = titleOffsets[targetIdx].id;
  if (newHash && location.hash.slice(1) !== newHash) {
    history.replaceState(null, "", `#${newHash}`);
  }

  lineIndex.value = titleOffsets[targetIdx].id;
};

onBeforeUnmount(() => {
  window.removeEventListener("scroll", anchorScroll);
  window.removeEventListener("scroll", subscribeTitleFocus);
});

// 点赞
const like = async () => {
  if (post.value?.is_liked) return;
  try {
    let likeRes: any = await likePost(id);
    let res: IBaseResponse = likeRes.data.value;
    if (res?.code === 0) {
      post.value!.is_liked = true;
      post.value!.like_count++;
    } else {
      console.log(res);
    }
  } catch (error) {
    console.log(error);
  }
};

const comment = ref();
const scrollToCommentArea = () => {
  window.scrollTo({
    top: comment.value.offsetTop - 60,
    behavior: "smooth",
  });
};

// 二维码
const qrcodeShow = ref(false);

import { useAlertStore } from "~/store/toast";
import {
  getComments,
  type IComment,
  type ICommentReplyRequest,
  type ICommentRequest,
  submitComment,
  submitCommentReply,
} from "~/api/comment";
import type { IPayInfo } from "~/api/config";
import { useConfigStore } from "~/store/config";

const toast = useAlertStore();

const copyLink = async () => {
  await navigator.clipboard.writeText(link.value);
  toast.showToast("复制成功！", 2000);
};

const { data: comments } = await useAsyncData(
  `post-${id}-comments`,
  async () => {
    try {
      let commentRes: any = await getComments(id);
      let res: IResponse<IPageData<IComment>> = commentRes.data.value;
      if (res && res.data) {
        if (res.code !== 0) {
          toast.showToast(res.message, 2000);
          return [];
        }
        return res.data?.list || [];
      }
    } catch (error: any) {
      toast.showToast(error.toString(), 2000);
      return [];
    }
  },
);

const submit = async (req: ICommentRequest) => {
  try {
    req.postId = id;
    let commentRes: any = await submitComment(req);
    if (commentRes.data.value === null) {
      if (commentRes.error.value.statusCode == 403) {
        toast.showToast("评论模块暂未开放！", 2000);
      } else {
        toast.showToast(commentRes.error.value.statusMessage, 2000);
      }
      return;
    }
    let res: IBaseResponse = commentRes.data.value;
    if (res.code !== 0) {
      toast.showToast(res.message, 2000);
      return;
    }
    toast.showToast("提交评论成功，待站长审核通过后将会通过邮件告知。", 3000);
    clearCommentReq();
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
};

const submitReply = async (req: ICommentReplyRequest, commentId: string) => {
  try {
    req.postId = id;
    let commentRes: any = await submitCommentReply(commentId, req);
    if (commentRes.data.value === null) {
      if (commentRes.error.value.statusCode == 403) {
        toast.showToast("评论模块暂未开放！", 2000);
      } else {
        toast.showToast(commentRes.error.value.statusMessage, 2000);
      }
      return;
    }
    let res: IBaseResponse = commentRes.data.value;
    if (res.code !== 0) {
      toast.showToast(res.message, 2000);
      return;
    }
    toast.showToast("提交评论成功，待站长审核通过后将会通过邮件告知。", 3000);
    clearCommentReplyReq();
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
};

const commentPost = ref();
const clearCommentReq = () => {
  if (commentPost.value) {
    commentPost.value.clearReq();
  }
};
const clearCommentReplyReq = () => {
  if (commentPost.value) {
    commentPost.value.clearReplyReq();
  }
};

const submitReply2Reply = async (
  req: ICommentReplyRequest,
  commentId: string,
) => {
  try {
    req.postId = id;
    let commentRes: any = await submitCommentReply(commentId, req);
    if (commentRes.data.value === null) {
      if (commentRes.error.value.statusCode == 403) {
        toast.showToast("评论模块暂未开放！", 2000);
      } else {
        toast.showToast(commentRes.error.value.statusMessage, 2000);
      }
      return;
    }
    let res: IBaseResponse = commentRes.data.value;
    if (res.code !== 0) {
      toast.showToast(res.message, 2000);
      return;
    }
    toast.showToast("提交评论成功，待站长审核通过后将会通过邮件告知。", 3000);
    clearReply2ReplyReq();
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
};

let description = post.value?.summary;
if (post.value?.meta_description) {
  description = post.value?.meta_description;
}
const clearReply2ReplyReq = () => {
  if (commentPost.value) {
    commentPost.value.clearReply2ReplyReq();
  }
};

useHead({
  title: `${post.value?.title} - ${
    configStore.seo_meta_config.title === ""
      ? configStore.website_info.website_name
      : configStore.seo_meta_config.title
  }`,
  meta: [{ name: "description", content: description }],
});
useSeoMeta({
  ogTitle: `${post.value?.title} - ${
    configStore.seo_meta_config.og_title === ""
      ? configStore.website_info.website_name
      : configStore.seo_meta_config.og_title
  }`,
  ogDescription: description,
  ogImage: configStore.seo_meta_config.og_image,
  twitterCard: "summary",
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

.share::before {
  content: "";
  position: absolute;
  left: -20px;
  top: 25%;
  border-width: 10px;
  border-style: solid;
  border-color: transparent #b7bbc4 transparent transparent;
  transform: translateY(-50%);
}

.qrcodeDiv::before {
  content: "";
  position: absolute;
  left: -20px;
  top: 25%;
  border-width: 10px;
  border-style: solid;
  border-color: transparent #b7bbc4 transparent transparent;
  transform: translateY(-50%);
}

@keyframes tilt {
  0%,
  100% {
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
  content: "";
  position: absolute;
  left: -20px;
  top: 10%;
  border-width: 10px;
  border-style: solid;
  border-color: transparent #b7bbc4 transparent transparent;
  transform: translateY(-50%);
}

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
