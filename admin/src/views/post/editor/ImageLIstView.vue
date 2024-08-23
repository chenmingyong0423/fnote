<template>
  <a-card>
    <div class="flex h-500px text-#8491a5">
      <div class="w-30% h-500px border-r-1 border-r-solid border-r-#f0f0f0 relative">
        <div>
          <ul class="h-full overflow-y-auto p-1">
            <li
              v-for="item in folders"
              :key="item.id"
              class=""
              :class="[
                'box-border flex items-center justify-between p-4 rounded-lg cursor-pointer hover:bg-gray-100 group h-54px m-y-2',
                {'bg-gray-100': state.selectedMenuItem == item.id}
              ]"
              @click="menuItemChanged(item.id)"
            >
              <div>
                <span class="text-gray-600">
                  <component :is="item.icon" />
                </span>
                <span class="ml-2 ">{{ item.name }}</span>
              </div>
              <div class="hidden group-hover:block">
                <a-button shape="circle" size="small" v-if="item.support_edit">
                  <template #icon>
                    <FormOutlined style="color: #8491a5;" />
                  </template>
                </a-button>

                <a-button shape="circle" size="small" v-if="item.support_delete">
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
            v-for="(image) in images"
            :key="image.id"
            @click="selectImage(image.id)"
            @mouseover="state.hoveredImgIndex = image.id"
            @mouseleave="state.hoveredImgIndex = null"
            :class="['relative box-border border rounded w-27 h-27 cursor-pointer',
                { 'border-blue-500': isSelected(image.id)
                }]"
          >
            <img :src="image.src" class="box-border w-27 h-27 object-cover" alt="image">

            <div v-if="isSelected(image.id)" class="w-full h-full absolute top-0">
              <div class="box-border bg-black opacity-20 w-full h-full"></div>
              <div class="absolute inset-0 border-2 border-solid border-blue-500"></div>
              <div class="absolute top-2 right-2">
                <span class="bg-blue-500 text-white text-xl rounded-full w-6 h-6 lh-6 flex items-center justify-center">√</span>
              </div>
            </div>

            <div v-if="state.hoveredImgIndex == image.id && !isSelected(image.id)" class="w-full h-full absolute top-0">
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
          <a-button type="primary" :disabled="state.selectedImgIndex==null" @click="insert">
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
import { type AssetFolderVO, GetAssetFolderList } from '@/interfaces/PostAsset'

const state = reactive({
  selectedMenuItem: '',
  selectedImgIndex: null,
  hoveredImgIndex: null
})

const folders = ref<AssetFolderVO[]>([])

const getFolders = async () => {
  try {
    const response = await GetAssetFolderList("image", "post-editor")
    folders.value = response.data.data?.list || []
    if (folders.value.length > 0) {
      folders.value.forEach(item => { item.icon = h(PictureOutlined) })
      state.selectedMenuItem = folders.value[0].id
    }
  } catch (error) {
    console.log(error)
  }
}

getFolders()

const images = [
  { id: "1", src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' },
  { id: "2", src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' },
  { id: "3", src: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg' }
]

const selectImage = (id) => {
  state.selectedImgIndex = state.selectedImgIndex === id ? null : id
}

const isSelected = (id) => {
  return state.selectedImgIndex === id
}

const emit = defineEmits(['insertImg'])

const insert = () => {
  // 从 images 里面找到对应元素
  const image = images.find((item) => item.id == state.selectedImgIndex)
  // 告诉父组件
  emit('insertImg', `![](${image.src})`)
}

const menuItemChanged = (id: string) => {
  state.selectedMenuItem = id
}
</script>