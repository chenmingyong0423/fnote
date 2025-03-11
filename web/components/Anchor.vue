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
      v-for="(anchor, index) in anchors"
      :style="{
        padding: `10px 0 10px ${anchor.depth * anchor.depth * anchor.depth}px`,
      }"
      :key="index"
    >
      <a
        style="cursor: pointer"
        class="p-l-6"
        :class="{
          'anchor_border text-5 text-#1e80ff font-bold active':
            anchor.id == lineIdx,
        }"
        :href="'#' + anchor.id"
        @click="click(anchor)"
        >{{ anchor.text }}</a
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  toc: Toc[];
  lineIndex: string;
}>();

const lineIdx = ref("");

type Toc = {
  id: string;
  depth: number;
  text: string;
  children: Toc[];
};

// 将 toc 逐渐展开并合并到 anchors 中
const anchors = ref<Toc[]>([]);
watch(
  () => props.toc,
  (newValue) => {
    anchors.value = [];
    expandToc(newValue, 1);
  },
  { immediate: true },
);

function expandToc(toc: Toc[], depth: number = 1) {
  toc.forEach((item) => {
    anchors.value.push(item);
    if (item.children) {
      expandToc(item.children, depth + 1);
    }
  });
}

const container = ref();

const click = (toc: Toc) => {
  lineIdx.value = toc.id;
};

watch(
  () => props.lineIndex,
  async (newValue, oldValue) => {
    if (newValue !== oldValue) {
      await nextTick();
      lineIdx.value = newValue;
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
