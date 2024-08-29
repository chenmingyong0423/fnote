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
                { 'bg-gray-100': state.selectedMenuItem == item.id }
              ]"
              @click="menuItemChanged(item.id)"
            >
              <div>
                <span class="text-gray-600">
                  <component :is="item.icon" />
                </span>
                <span class="ml-2">{{ item.name }}</span>
              </div>
              <div class="hidden group-hover:block">
                <a-button
                  shape="circle"
                  size="small"
                  v-if="item.support_edit"
                  @click="preEdit(item)"
                >
                  <template #icon>
                    <FormOutlined style="color: #8491a5" />
                  </template>
                </a-button>

                <a-button shape="circle" size="small" v-if="item.support_delete">
                  <template #icon>
                    <a-popconfirm
                      title="确认删除？"
                      ok-text="是"
                      cancel-text="否"
                      @confirm="deleteAssetFolder(item.id)"
                    >
                      <DeleteOutlined style="color: #8491a5" />
                    </a-popconfirm>
                  </template>
                </a-button>
              </div>
            </li>
          </ul>
        </div>
        <a-button
          class="absolute bottom-0 text-#8491a5"
          @click="
            () => {
              visible = true
              modalLabel = '新增分类'
            }
          "
        >
          <template #icon>
            <FolderAddOutlined />
          </template>
          新增分类
        </a-button>
        <a-modal v-model:visible="visible" :title="modalLabel" @ok="handleOk" @cancel="cancel">
          <a-input addon-before="分类名称" v-model:value="assetFolder.name" />
        </a-modal>
      </div>
      <div class="w-69% ml-1% flex flex-col gap-y-2 relative">
        <div class="flex gap-x-1">
          <SimpleUpload
            @success:imageUrl="uploadAsset"
            :authorization="userStore.token"
            :action="serverHost + '/admin-api/files/upload'"
            label="上传图片"
            :fileTypes="['image/jpeg', 'image/png']"
            :maxSize="1048576"
          />
          <a-button @click="visible4ExistImageModal = true">
            <template #icon><FileImageOutlined /></template>
            选择已有图片
          </a-button>
        </div>
        <div class="flex gap-2 flex-wrap">
          <PreviewImg v-model="visible4PreviewImgModal" :image-url="previewImgUrl" />
          <div
            v-for="image in images"
            :key="image.id"
            @click="selectImage(image.id)"
            class="relative box-border border rounded w-27 h-27 cursor-pointer"
          >
            <img
              :src="serverHost + image.content"
              class="box-border w-27 h-27 object-contain"
              alt="image"
            />

            <div
              :class="[
                'w-full h-full absolute top-0 duration-300 ease-in-out group',
                {
                  'bg-black bg-opacity-40': isSelected(image.id),
                  'hover:bg-black hover:bg-opacity-40': !isSelected(image.id)
                }
              ]"
            >
              <div
                class="hidden absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 group-hover:block duration-300 ease-in-out group"
              >
                <div class="flex space-x-2">
                  <a-button shape="circle" :ghost="true" @click.stop="preview(image.content)">
                    <template #icon><EyeOutlined /></template>
                  </a-button>
                </div>
              </div>
            </div>

            <div v-if="isSelected(image.id)" class="absolute top-2 right-2 opacity-100">
              <span
                class="bg-blue-500 text-white text-xl rounded-full w-6 h-6 lh-6 flex items-center justify-center"
                >√</span
              >
            </div>
          </div>
        </div>
        <div class="absolute bottom-0 w-full flex gap-x-1 justify-end">
          <a-popconfirm title="确认删除？" ok-text="是" cancel-text="否" @confirm="deleteAsset">
            <a-button type="primary" :disabled="state.selectedImgIndex == ''"> 删除图片 </a-button>
          </a-popconfirm>

          <a-button type="primary" :disabled="state.selectedImgIndex == ''" @click="insert">
            插入图片
          </a-button>
        </div>
      </div>
    </div>
    <ImageList v-model="visible4ExistImageModal" @insertImg="uploadAsset"></ImageList>
  </a-card>
</template>

<script setup lang="ts">
import { defineEmits, h, reactive, ref, watch } from 'vue'
import {
  PictureOutlined,
  FolderAddOutlined,
  FormOutlined,
  DeleteOutlined,
  FileImageOutlined,
  EyeOutlined
} from '@ant-design/icons-vue'
import {
  AddAsset,
  AddAssetFolder,
  type AddAssetFolderRequest,
  type AssetFolderVO,
  type AssetRequest,
  type AssetVO,
  DeleteAsset,
  DeleteAssetFolder,
  EditAssetFolderName,
  GetAssetFolderList,
  GetAssetList
} from '@/interfaces/Asset'
import { message } from 'ant-design-vue'
import originalAxios from 'axios'
import { useUserStore } from '@/stores/user'
import SimpleUpload from '@/components/upload/SimpleUpload.vue'
import ImageList from '@/components/file/ImageList.vue'
import PreviewImg from '@/components/image/PreviewImg.vue'

const state = reactive({
  selectedMenuItem: '',
  selectedImgIndex: ''
})

const folders = ref<AssetFolderVO[]>([])

const getFolders = async () => {
  try {
    const response = await GetAssetFolderList('image', 'post-editor')
    folders.value = response.data.data?.list || []
    if (folders.value.length > 0) {
      folders.value.forEach((item) => {
        item.icon = h(PictureOutlined)
      })
      state.selectedMenuItem = folders.value[0].id
    }
  } catch (error) {
    console.log(error)
  }
}

getFolders()

const images = ref<AssetVO[]>([])

const selectImage = (id: string) => {
  state.selectedImgIndex = state.selectedImgIndex == id ? '' : id
}

const isSelected = (id: string) => {
  return state.selectedImgIndex == id
}

const emit = defineEmits(['insertImg'])

const insert = () => {
  // 从 images 里面找到对应元素
  const image = images.value.find((item) => item.id == state.selectedImgIndex)
  // 告诉父组件
  emit('insertImg', `![](${image?.content})`)
}

const menuItemChanged = (id: string) => {
  state.selectedMenuItem = id
}

const visible = ref(false)
const modalLabel = ref('')

const assetFolder = reactive<AddAssetFolderRequest>({
  name: '',
  asset_type: 'image',
  type: 'post-editor',
  support_delete: true,
  support_edit: true,
  support_add: true
} as AddAssetFolderRequest)

const handleOk = () => {
  if (modalLabel.value === '新增分类') {
    addAssetFolder()
  } else {
    editAssetFolder()
  }
}

const addAssetFolder = async () => {
  if (!assetFolder.name) {
    message.error('请输入分类名称')
  } else {
    try {
      const response: any = await AddAssetFolder(assetFolder)
      if (response.data.code !== 0) {
        message.error(response.data.message)
        return
      }
      message.success('添加成功')
      await getFolders()
      assetFolder.name = ''
      visible.value = false
    } catch (error) {
      console.log(error)
      if (originalAxios.isAxiosError(error)) {
        // 这是一个由 axios 抛出的错误
        if (error.response) {
          if (error.response.status === 409) {
            message.error('分类名称重复')
            return
          }
        } else if (error.request) {
          // 请求已发出，但没有收到响应
          console.log('No response received:', error.request)
        } else {
          // 在设置请求时触发了一个错误
          console.log('Error Message:', error.message)
        }
      }
      message.error('添加失败')
    }
  }
}

const preEdit = (assetFolderVO: AssetFolderVO) => {
  modalLabel.value = '编辑分类'
  visible.value = true
  assetFolder.id = assetFolderVO.id
  assetFolder.name = assetFolderVO.name
}

const resetAssetFolder = () => {
  assetFolder.id = ''
  assetFolder.name = ''
  assetFolder.asset_type = 'image'
  assetFolder.type = 'post-editor'
  assetFolder.support_delete = true
  assetFolder.support_edit = true
  assetFolder.support_add = true
}

const editAssetFolder = async () => {
  if (!assetFolder.name) {
    message.error('请输入分类名称')
  } else {
    try {
      const response: any = await EditAssetFolderName(assetFolder.id, assetFolder.name)
      if (response.data.code !== 0) {
        message.error(response.data.message)
        return
      }
      message.success('编辑成功')
      await getFolders()
      resetAssetFolder()
      visible.value = false
    } catch (error) {
      console.log(error)
      if (originalAxios.isAxiosError(error)) {
        // 这是一个由 axios 抛出的错误
        if (error.response) {
          if (error.response.status === 409) {
            message.error('分类名称重复')
            return
          }
        } else if (error.request) {
          // 请求已发出，但没有收到响应
          console.log('No response received:', error.request)
        } else {
          // 在设置请求时触发了一个错误
          console.log('Error Message:', error.message)
        }
      }
      message.error('编辑失败')
    }
  }
}

const cancel = () => {
  resetAssetFolder()
}

const deleteAssetFolder = async (id: string) => {
  try {
    const response: any = await DeleteAssetFolder(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getFolders()
  } catch (error) {
    console.log(error)
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.status === 404) {
          message.error('分类不存在')
          return
        }
      } else if (error.request) {
        // 请求已发出，但没有收到响应
        console.log('No response received:', error.request)
      } else {
        // 在设置请求时触发了一个错误
        console.log('Error Message:', error.message)
      }
    }
    message.error('删除失败')
  }
}

const getImages = async () => {
  try {
    const response: any = await GetAssetList(state.selectedMenuItem)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    images.value = response.data.data?.list || []
  } catch (error) {
    console.log(error)
    message.error('获取图片列表失败')
  }
}

watch(
  () => state.selectedMenuItem,
  () => {
    getImages()
  }
)

const serverHost = import.meta.env.VITE_API_HOST
const userStore = useUserStore()

const uploadAsset = async (fileId: string, fileUrl: string) => {
  try {
    const response: any = await AddAsset(state.selectedMenuItem, {
      content: fileUrl,
      asset_type: 'image',
      type: 'post-editor',
      metadata: {
        file_id: fileId
      }
    } as AssetRequest)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('上传成功')
    await getImages()
  } catch (error) {
    console.log(error)
    message.error('上传失败')
  }
}

const visible4ExistImageModal = ref(false)

const visible4PreviewImgModal = ref(false)
const previewImgUrl = ref('')
const preview = (imgUrl: string) => {
  previewImgUrl.value = serverHost + imgUrl
  visible4PreviewImgModal.value = true
}

const deleteAsset = async () => {
  try {
    if (state.selectedMenuItem === '' || state.selectedImgIndex === '') {
      message.error('请选择图片')
      return
    }
    const response: any = await DeleteAsset(state.selectedMenuItem, state.selectedImgIndex)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    state.selectedImgIndex = ''
    await getImages()
  } catch (error) {
    console.log(error)
    message.error('删除失败')
  }
}
</script>
