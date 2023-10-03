<template>
    <div class="bg-#fff rounded-15 p20 dark_bg_gray dark_text_white ">
        <div class="text-24">目录</div>
        <el-divider></el-divider>
        <el-scrollbar height="400px">
            <div v-for="anchor, index in titleList" :style="{ 'padding-left': `${anchor.indent * 20 + 20}px` }"
                @click="rollTo(anchor, index)" :class="index === heightTitle ? 'title-active' : ''"
                class="cursor-pointer text-16 py10 my10">
                <a>{{ anchor.title }}</a>
            </div>
        </el-scrollbar>
    </div>
</template>

<script setup>
import { throttle } from '~/utils'
const props = defineProps(['text'])
// markdown-文章标题列表
const titleList = ref([])

// markdown-生成标题
function getTitle() {

    // 使用js选择器，获取对应的h标签，组合成列表
    const anchors = document.querySelectorAll(
        '.v-md-editor-preview h1,h2,h3,h4,h5,h6'
    )
    console.log('获取到标题列表', anchors);
    // 删除标题头尾的空格
    const titles = Array.from(anchors).filter((title) => !!title.innerText.trim());
    console.log('处理完成，清理空标题', titles);
    // 当文章h标签为空时，直接返回
    if (!titles.length) {
        titleList.value = [];
        return;
    }
    // 从h标签属性中，提取相关信息
    const hTags = Array.from(new Set(titles.map((title) => title.tagName))).sort();
    console.log('包含以下标签', hTags);
    titleList.value = titles.map((item) => ({
        title: item.innerText, // 标题内容
        lineIndex: item.getAttribute('data-v-md-line'), // 标签line id
        indent: hTags.indexOf(item.tagName), // 标签层级
        height: item.offsetTop, // 标签距离顶部距离
    }));
}
// markdown-当前高亮的标题index
const heightTitle = ref(0)
// markdown-标题跳转
const rollTo = (anchor, index) => {
    // 获取要跳转的标签的lineIndex
    const { lineIndex } = anchor;
    // 查找lineIndex对应的元素对象
    const heading = document.querySelector(
        `.v-md-editor-preview [data-v-md-line="${lineIndex}"]`
    );
    // 页面跳转
    if (heading) {
        scrollTo({
            top: heading.offsetTop,
            behavior: 'smooth'
        })
    }
    // 修改当前高亮的标题
    heightTitle.value = index
}

// markdown-页面滚动。
const scroll = () => {
    // 监听屏幕滚动时防抖（在规定的时间内触发的事件，只执行最后一次，降低性能开销）
    return () => {
        throttle(jump, 100)
    }
}
const jump = () => {
    let scrollTop = window.pageYOffset
    // console.log(window.pageYOffset)
    const absList = [] // 各个h标签与当前距离绝对值
    titleList.value.forEach((item) => {
        absList.push(Math.abs(item.height - scrollTop))
        // console.log('height', item.height);
    })
    // 屏幕滚动距离与标题高度最近的index高亮
    heightTitle.value = absList.indexOf(Math.min.apply(null, absList))
}
onMounted(() => {
    // 生成文章标题列表
    getTitle()
    window.addEventListener('scroll', throttle(jump, 100))
})
onBeforeUnmount(() => {
    window.removeEventListener('scroll', throttle(jump, 100))
})
</script>
<style scoped>
.title-active {
    color: cornflowerblue;
    border-left: 5px solid cornflowerblue;

}
</style>
