<template>
  <a-card>
    <div class="flex h-500px text-#8491a5">
      <div class="w-30% h-500px border-r-1 border-r-solid border-r-#f0f0f0 relative">
        <div>
          <ul class="h-full overflow-y-auto">
            <li
              v-for="item in items"
              :key="item.key"
              class="box-border flex items-center justify-between p-4 rounded-lg cursor-pointer hover:bg-blue-100 group h-54px"
            >
              <div>
                <span class="text-gray-600">
                  <component :is="item.icon" />
                </span>
                <span class="ml-2 ">{{ item.title }}</span>
              </div>
              <div class="hidden group-hover:block">
                <a-button shape="circle" size="small">
                  <template #icon>
                    <FormOutlined style="color: #8491a5;" />
                  </template>
                </a-button>

                <a-button shape="circle" size="small">
                  <template #icon>
                    <DeleteOutlined style="color: #8491a5;" />
                  </template>
                </a-button>

              </div>
            </li>
          </ul>
        </div>
        <a-button class="absolute bottom-0 text-#8491a5">
          <template #icon>
            <FolderAddOutlined />
          </template>
          新增分类
        </a-button>
      </div>
      <div class="w-69% ml-1% flex flex-col gap-y-2 relative">
        <div>
          <a-button class="text-#8491a5">
            <template #icon>
              <PictureOutlined />
            </template>
            上传图片
          </a-button>
        </div>
        <div class="flex gap-2">
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
                <span
                  class="rounded-full border-2 border-solid border-white  w-6 h-6 flex items-center justify-center"></span>
              </div>
            </div>
          </div>
        </div>
        <div class="absolute bottom-0 w-full text-right">
          <a-button type="primary" :disabled="selectedIndex==null" @click="insert">
            插入图片
          </a-button>
        </div>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { defineEmits, h, reactive, ref } from 'vue'
import { PictureOutlined, FolderAddOutlined, FormOutlined, DeleteOutlined } from '@ant-design/icons-vue'

const state = reactive({
  collapsed: false,
  selectedKeys: ['none'],
  openKeys: ['none']
})

const items = reactive([
  {
    key: 'none',
    icon: () => h(PictureOutlined),
    title: '未分类'
  },
  {
    key: 'post-fixed-material',
    icon: () => h(PictureOutlined),
    label: '文章固定素材',
    title: '文章固定素材'
  }
])

const images = [
  { src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' },
  { src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' },
  { src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' }
]

const selectedIndex = ref(null)
const hoveredIndex = ref(null)

const selectImage = (index) => {
  selectedIndex.value = selectedIndex.value === index ? null : index
}

const isSelected = (index) => {
  return selectedIndex.value === index
}

const emit = defineEmits(['insertImg'])

const insert = () => {
  // 告诉父组件
  emit('insertImg', '陈明勇')
}
</script>