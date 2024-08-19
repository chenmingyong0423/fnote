<template>
  <a-card>
    <div class="flex min-h-500px">
      <div class="w-30%">
        <a-menu
          v-model:openKeys="state.openKeys"
          v-model:selectedKeys="state.selectedKeys"
          mode="inline"
          theme="light"
          :items="items"
          @click="itemClick"
        ></a-menu>
      </div>
      <div class="w-69% ml-1% flex gap-2">
        <div
          v-for="(image, index) in images"
          :key="index"
          @click="selectImage(index)"
          @mouseover="hoveredIndex = index"
          @mouseleave="hoveredIndex = null"
          :class="['relative box-border border rounded w-27 h-27 cursor-pointer',
                { 'border-blue-500': isSelected(index)
                }]"
        >
          <img :src="image.src" class="box-border w-27 h-27 object-cover" alt="image">

          <div v-if="isSelected(index)" class="w-full h-full absolute top-0">
            <div class="box-border bg-black opacity-20 w-full h-full"></div>
            <div class="absolute inset-0 border-2 border-solid border-blue-500"></div>
            <div class="absolute top-2 right-2">
              <span class="bg-blue-500 text-white text-xl rounded-full w-6 h-6 lh-6 flex items-center justify-center">√</span>
            </div>
          </div>

          <div v-if="hoveredIndex == index && !isSelected(index)" class="w-full h-full absolute top-0">
            <div class="box-border bg-black opacity-20 w-full h-full"></div>
            <div class="absolute inset-0 border-2 border-solid border-blue-500"></div>
            <div class="absolute top-2 right-2">
              <span class="rounded-full border-2 border-solid border-white  w-6 h-6 flex items-center justify-center"></span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { h, reactive, ref } from 'vue'
import { PictureOutlined } from '@ant-design/icons-vue'

const state = reactive({
  collapsed: false,
  selectedKeys: ['none'],
  openKeys: ['none']
})

const items = reactive([
  {
    key: 'none',
    icon: () => h(PictureOutlined),
    label: '未分类',
    title: '未分类'
  },
  {
    key: 'post-fixed-material',
    icon: () => h(PictureOutlined),
    label: '文章固定素材',
    title: '文章固定素材'
  }
])

const itemClick = (item: any) => {
  console.log(item)
}

const images = [
  { src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' },
  { src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' },
  { src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' }
];

const selectedIndex = ref(null);
const hoveredIndex = ref(null);

const selectImage = (index) => {
  selectedIndex.value = selectedIndex.value === index ? null : index;
};

const isSelected = (index) => {
  return selectedIndex.value === index;
};
</script>