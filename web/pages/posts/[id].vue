<template>
  <div class="flex w-full">
    <div class="mt-10 w-10%">
      <div class="flex flex-col gap-y-3 items-center fixed">
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
          <span
              class="i-ph:thumbs-up w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
        <span
            class="i-ph-chats-duotone w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
        <span
            class="i-ph:share-fat-light w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400"></span>
        </div>
        <div
            class="group flex items-center justify-center w-12 h-12 border-rounded-50% bg-white p-2 cursor-pointer hover-bg-#1e80ff duration-200">
           <span
               class="w-8 h-8 text-gray group-hover:scale-120 group-hover:text-white duration-400 text-5 text-center">赏</span>
        </div>
      </div>
    </div>
    <div class="w-58% ml-1% mr-1%">
      <div class="bg-white mb-5 b-rounded-4">
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
          <v-md-preview :text="post?.content" @copy-code-success="handleCopyCodeSuccess"
                        class="dark_text_white lt-lg:important:p0" @change="generateAnchors"></v-md-preview>
        </div>
      </div>
      <!-- 版权声明 -->
      <div class="copyright b-rounded-4 bg-white p-8">
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
    </div>
    <div class="flex flex-col w-30%">
      <Profile class="mb-5"></Profile>
      <div ref="anchor">
        <Anchor :htmlContent="htmlContent" :lineIndex="lineIndex" @handleAnchorClick="handleAnchorClick"></Anchor>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {getPostsById} from "~/api/post";
import type {IPostDetail} from "~/api/post";
import type {IResponse} from "~/api/http";
import {onMounted, ref} from "vue";
import {useHomeStore} from '~/store/home';

const info = useHomeStore()
const domain = info.master_info.domain;
const route = useRoute()
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
  console.log(isScrolling.value)
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
</script>