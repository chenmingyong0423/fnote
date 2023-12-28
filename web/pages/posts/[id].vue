<template>
  <div class="flex w-full">
    <div class="mt-10 w-5%">
      <div class="flex flex-col gap-y-3 items-center fixed ">
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white  p-2 cursor-pointer duration-200 dark:text-dtc dark_bg_gray relative"
            :class="{' hover:bg-#1e80ff': !post?.is_liked}" @click="like">
          <span
              class="i-ph:thumbs-up w-8 h-8 duration-400"
              :class="{'group-hover:scale-120 group-hover:text-white text-gray': !post?.is_liked, 'text-#1e80ff': post?.is_liked}"></span>
          <span
              class="absolute translate-x-11/10 -translate-y-11/10 -translate-y-1/2 bg-#1e80ff text-white text-xs rounded-full w-6 h-6 flex items-center justify-center">
              {{ post?.like_count }}
          </span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray relative"
            @click="scrollToCommentArea">
          <span
              class="i-ph-chats-duotone w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
          <span
              class="absolute transform translate-x-11/10 -translate-y-11/10 bg-#1e80ff text-white text-xs rounded-full w-6 h-6 flex items-center justify-center">
              {{ post?.comment_count }}
          </span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray relative">
          <span
              class="i-ph:share-fat-light w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
          <div class="share w-80% hidden absolute transform translate-x-8/6 -translate-y--1/4 group-hover:block">
            <div
                class=" flex flex-col gap-y-3  bg-white text-white text-xs b-rounded-4 w-80% items-center justify-center p-x-1 p-y-2">
              <span class="i-bi:wechat w-6 h-6 text-green"></span>
              <hr class="w-50% border-gray-1 b-rounded-4">
              <a class="i-bi:tencent-qq w-6 h-6 text-black"
                 :href="`https://connect.qq.com/widget/shareqq/index.html?url=${domain}${path}&title=${post?.title}&pics=${post?.cover_img}`" target="_blank"></a>
              <hr class="w-50% border-gray-1 b-rounded-4">
              <a class="i-bi:sina-weibo w-6 h-6 text-red" :href="`https://service.weibo.com/share/share.php?sharesource=weibo&title=${post?.title}，原文链接：${domain}${path}&pic=${post?.cover_img}`" target="_blank"></a>
              <hr class="w-50% border-gray-1 b-rounded-4">
              <span class="i-bi:link-45deg w-6 h-6 text-black"></span>
            </div>
          </div>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200 dark:text-dtc dark_bg_gray">
           <span
               class="w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400 text-5 text-center">赏</span>
        </div>
        <div>
          <div ref="qrCodeEl">

          </div>
          <p @click="generateQRCode">扫一扫</p>
        </div>
      </div>
    </div>
    <div class="w-63% ml-1% mr-1%">
      <div class="bg-white mb-5 b-rounded-4 dark:text-dtc dark_bg_gray">
        <!--  文章标题  -->
        <div class="text-10 font-bold text-center p-1">{{ post?.title }}</div>
        <!--  文章 meta  -->
        <div class="flex items-center gap-x-2 text-4 justify-center p-1 text-gray-4">
          <div>{{ post?.author }}</div>
          <div>{{ $dayjs(post?.create_time).format('YYYY-MM-DD HH:mm:ss') }}</div>
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
      <!-- 版权声明 -->
      <div class="copyright b-rounded-4 bg-white p-8 dark:text-dtc dark_bg_gray">
        <p class="mb-5"><span style="color: rgb(14, 136, 235);font-weight: bold;">本文链接：</span><a
            class="text-#00bd7e hover:bg-#00bd7e33"
            :href="`${domain}/${id}`" target="_blank">{{ `${domain}/posts/${id}` }}</a></p>
        <p><span style="color: rgb(14, 136, 235);font-weight: bold;">版权声明：</span>本文由 <span
            style="color: rgb(14, 136, 235);">{{ post?.author }}</span> 原创发布，如需转载请遵循 <a
            class="text-#00bd7e hover:bg-#00bd7e33"
            href="https://creativecommons.org/licenses/by-nc-sa/4.0/deed.zh" target="_blank">署名-非商业性使用-相同方式共享
          4.0 国际
          (CC BY-NC-SA 4.0)</a> 许可协议授权</p>
      </div>
      <!-- 评论区 -->
      <div ref="comment">
        评论区
      </div>
    </div>
    <div class="flex flex-col w-30%">
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
import type {IResponse, IBaseResponse} from "~/api/http";
import {onMounted, ref} from "vue";
import {useHomeStore} from '~/store/home';
import VMdPreview from "@kangc/v-md-editor/lib/preview";

const info = useHomeStore()

const isBlackMode = computed(() => info.isBlackMode)

const domain = info.master_info.domain;
const route = useRoute()
const path: string = route.path
const id: string = String(route.params.id)
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
  titles.forEach((title, index) => {
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
}

// 二维码
const qrCodeEl = ref();

onMounted(() => {
  if(process.client) {
    // if (qrCodeEl.value) {
    //   console.log(1234)
    //   import('qrcodejs2-fixes').then(QRCode => {
    //     new QRCode(qrCodeEl.value, {
    //       text: `${domain}/posts/${id}`,
    //       width: 100,
    //       height: 100,
    //       colorDark: "#000000",
    //       colorLight: "#ffffff",
    //       correctLevel: QRCode.default.CorrectLevel.H
    //     });
    //   })
    // }
  }
})

const generateQRCode = () => {
  console.log(1)
  // if (process.client) {
  //   console.log(2)
  //   console.log(qrCodeEl.value)
  //   if (qrCodeEl.value) {
  //     console.log(3)
  //     import('qrcodejs2-fixes').then(QRCode => {
  //       new QRCode(qrCodeEl.value, {
  //         text: `${domain}/posts/${id}`,
  //         width: 100,
  //         height: 100,
  //         colorDark: "#000000",
  //         colorLight: "#ffffff",
  //         correctLevel: QRCode.CorrectLevel.H
  //       });
  //     })
  //   }
  // }
}

</script>

<style scoped>
.dark /deep/ .v-md-pre-wrapper {
  background-color: rgba(10, 0, 0, 0.1) !important;
}

.dark /deep/ .v-md-pre-wrapper {
  background-color: rgba(10, 0, 0, 0.1) !important;
}

.dark /deep/ code {
  color: white !important;
}

.dark /deep/ .line-numbers-mode:after {
  background-color: rgba(10, 0, 0, 0.1) !important;
  border: 0 !important;
}

/* 根据需要定制不同代码语言或元素的样式 */
.dark /deep/ .hljs-keyword, .hljs-selector-tag, .hljs-literal {
  color: #ff7b72 !important;
}

.dark /deep/ .hljs-string {
  color: #a5d6ff !important;
}

.dark /deep/ .hljs-title {
  color: #a5d6ff !important;
}

.dark /deep/ .hljs-type {
  color: #cc880a !important;
}

.dark /deep/ .github-markdown-body table tr {
  background-color: rgba(10, 0, 0, 0.1) !important;

}

.dark /deep/ .github-markdown-body blockquote {
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

</style>
