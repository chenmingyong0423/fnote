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
        <div>
          <SimpleUpload
            @success:imageUrl="uploadAsset"
            :authorization="userStore.token"
            :action="serverHost + '/admin-api/files/upload'"
            label="上传图片"
            :fileTypes="['image/jpeg', 'image/png']"
            :maxSize="1048576"
          />
        </div>
        <div class="flex gap-2">
          <div
            v-for="image in images"
            :key="image.id"
            @click="selectImage(image.id)"
            @mouseover="state.hoveredImgIndex = image.id"
            @mouseleave="state.hoveredImgIndex = ''"
            :class="[
              'relative box-border border rounded w-27 h-27 cursor-pointer',
              { 'border-blue-500': isSelected(image.id) }
            ]"
          >
            <img
              :src="serverHost + image.content"
              class="box-border w-27 h-27 object-cover"
              alt="image"
            />

            <div v-if="isSelected(image.id)" class="w-full h-full absolute top-0">
              <div class="box-border bg-black opacity-20 w-full h-full"></div>
              <div class="absolute inset-0 border-2 border-solid border-blue-500"></div>
              <div class="absolute top-2 right-2">
                <span
                  class="bg-blue-500 text-white text-xl rounded-full w-6 h-6 lh-6 flex items-center justify-center"
                  >√</span
                >
              </div>
            </div>

            <div
              v-if="state.hoveredImgIndex == image.id && !isSelected(image.id)"
              class="w-full h-full absolute top-0"
            >
              <div class="box-border bg-black opacity-20 w-full h-full"></div>
              <div class="absolute inset-0 border-2 border-solid border-blue-500"></div>
              <div class="absolute top-2 right-2">
                <span
                  class="rounded-full border-2 border-solid border-white w-6 h-6 flex items-center justify-center"
                ></span>
              </div>
            </div>
          </div>
        </div>
        <div class="absolute bottom-0 w-full text-right">
          <a-button type="primary" :disabled="state.selectedImgIndex == null" @click="insert">
            插入图片
          </a-button>
        </div>
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { defineEmits, h, reactive, ref, watch } from 'vue'
import {
  PictureOutlined,
  FolderAddOutlined,
  FormOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import {
  AddAsset,
  AddAssetFolder,
  type AddAssetFolderRequest,
  type AssetFolderVO, type AssetRequest,
  type AssetVO,
  DeleteAssetFolder,
  EditAssetFolderName,
  GetAssetFolderList,
  GetAssetList
} from '@/interfaces/PostAsset'
import { message } from 'ant-design-vue'
import originalAxios from 'axios'
import { useUserStore } from '@/stores/user'
import SimpleUpload from '@/components/upload/SimpleUpload.vue'

const state = reactive({
  selectedMenuItem: '',
  selectedImgIndex: '',
  hoveredImgIndex: ''
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

const uploadAsset = async (content: string) => {
  try {
    const response: any = await AddAsset(state.selectedMenuItem, {
      content: content,
      asset_type: 'image',
      type: 'post-editor'
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
</script>
