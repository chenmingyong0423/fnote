<template>
  <div
    ref="container"
    class="flex flex-col bg-white b-rounded-4 max-h-[700px] overflow-y-auto"
  >
    <div
      class="line-height-6 text-6 light_border_bottom p-b-5 dark_text_white p-2 pl-4 pt-5"
    >
      目录
    </div>
    <div
      v-for="(anchor, index) in titles"
      :style="{ padding: `10px 0 10px ${anchor.indent * 1}px` }"
      @click="click(anchor)"
      :key="index"
      class="cursor-pointer"
    >
      <a
        style="cursor: pointer"
        class="p-l-6"
        :class="{
          'anchor_border text-5 text-#1e80ff font-bold active':
            anchor.lineIndex == lineIndex,
        }"
        >{{ anchor.title }}</a
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  htmlContent: string;
  lineIndex: string;
}>();
const container = ref();

const emit = defineEmits(["handleAnchorClick"]);

const click = (anchor: Title) => {
  emit("handleAnchorClick", anchor.lineIndex);
};

type Title = {
  title: string;
  lineIndex: string | null;
  indent: number;
};

// Refs for the component and titles
const titles = ref<Title[]>([]);

const generateAnchors = (html: string) => {
  if (html != "") {
    // 创建一个新的 DOMParser 实例
    const parser = new DOMParser();
    // 解析字符串为 HTML 文档
    const doc = parser.parseFromString(html, "text/html");
    const anchors = doc.querySelectorAll("h1,h2,h3,h4,h5,h6");
    const extractedTitles = Array.from(anchors)
      .filter((title) => !!title.innerHTML?.trim())
      .map((el) => ({
        title: el.textContent ?? "",
        lineIndex: el.getAttribute("data-v-md-line") || "",
        indent: 0, // Initialize indent, will be set later
        tagName: el.tagName,
      }));
    if (!extractedTitles.length) {
      titles.value = [];
      return;
    }

    const hTags = Array.from(
      extractedTitles.map((title) => title.tagName),
    ).sort();
    titles.value = extractedTitles.map((el) => ({
      ...el,
      indent: hTags.indexOf(el.tagName),
    }));
  }
};

watch(
  () => props.htmlContent,
  (newValue) => {
    generateAnchors(newValue);
  },
  { immediate: true },
);

watch(
  () => props.lineIndex,
  async (newValue, oldValue) => {
    if (newValue !== oldValue) {
      await nextTick();
      const activeElement = container.value.querySelector(".active");
      if (activeElement) {
        const elementPosition =
          activeElement.offsetTop - container.value.offsetTop;
        // 判断元素是否超过当前 container 高度的一半
        if (elementPosition > container.value.clientHeight / 2) {
          container.value.scrollTop =
            elementPosition - container.value.clientHeight / 2;
        } else {
          container.value.scrollTop = 0;
        }
        // container.value.scrollTop = elementPosition - 20;
      }
    }
  },
);
</script>
