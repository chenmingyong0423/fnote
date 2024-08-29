<template>
  <a-modal
    v-model:visible="visible"
    title="图片列表"
    @ok="handleOk"
    @cancel="cancel"
    :destroyOnClose="true"
  >
    <div class="flex gap-2 flex-wrap">
      <PreviewImg v-model="visible4PreviewImgModal" :image-url="previewImgUrl" />
      <div
        v-for="image in images"
        :key="image.file_id"
        @click="selectImage(image.file_id)"
        class="relative box-border border rounded w-27 h-27 cursor-pointer"
      >
        <img
          :src="serverHost + image.url"
          class="box-border w-27 h-27 object-contain"
          alt="image"
        />

        <div
          :class="[
            'w-full h-full absolute top-0 duration-300 ease-in-out group',
            {
              'bg-black bg-opacity-40': isSelected(image.file_id),
              'hover:bg-black hover:bg-opacity-40': !isSelected(image.file_id)
            }
          ]"
        >
          <div
            class="hidden absolute z-100 opacity-100 top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 group-hover:block duration-300 ease-in-out group"
          >
            <div class="flex space-x-2">
              <a-button shape="circle" :ghost="true" @click.stop="preview(image.url)">
                <template #icon><EyeOutlined /></template>
              </a-button>
            </div>
          </div>
        </div>

        <div v-if="isSelected(image.file_id)" class="absolute top-2 right-2 opacity-100">
          <span
            class="bg-blue-500 text-white text-xl rounded-full w-6 h-6 lh-6 flex items-center justify-center"
            >√</span
          >
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { defineEmits, reactive, ref, watch } from 'vue'
import { type FileVO, GetFileList, type PageRequest } from '@/interfaces/File'
import { EyeOutlined } from '@ant-design/icons-vue'
import PreviewImg from '@/components/image/PreviewImg.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: true
  }
})

const visible = ref<Boolean>(props.modelValue)

watch(
  () => props.modelValue,
  (value: Boolean) => {
    visible.value = value
    if (value == true) {
      getImages()
    }
  }
)

const images = ref<FileVO[]>([])
const serverHost = import.meta.env.VITE_API_HOST
const state = reactive({
  selectedImgIndex: ''
})
const emit = defineEmits(['insertImg', 'update:modelValue'])

const selectImage = (id: string) => {
  state.selectedImgIndex = state.selectedImgIndex == id ? '' : id
}

const isSelected = (image: string) => {
  return state.selectedImgIndex == image
}

const handleOk = () => {
  // 告诉父组件
  const image = images.value.find((item) => item.file_id == state.selectedImgIndex)
  if (image) {
    state.selectedImgIndex = ''
    emit('insertImg', image.file_id, image.url)
    emit('update:modelValue', false)
    visible.value = false
  }
}

const cancel = () => {
  state.selectedImgIndex = ''
  emit('update:modelValue', false)
  visible.value = false
}

const getImages = async () => {
  try {
    const response = await GetFileList({
      pageNum: 1,
      pageSize: 10,
      fileType: ['image/png', 'image/jpeg']
    } as PageRequest)
    images.value = response.data.data?.list || []
  } catch (error) {
    console.log(error)
  }
}

if (visible.value) {
  getImages()
}

const visible4PreviewImgModal = ref(false)
const previewImgUrl = ref('')
const preview = (imgUrl: string) => {
  previewImgUrl.value = serverHost + imgUrl
  visible4PreviewImgModal.value = true
}
</script>
