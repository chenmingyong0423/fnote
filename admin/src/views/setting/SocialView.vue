<template>
  <div>
    <a-button type="primary" @click="visible = true">新增社交信息</a-button>
    <a-modal
      v-model:open="visible"
      title="新增分类"
      ok-text="提交"
      cancel-text="取消"
      @ok="addSocialInfo"
    >
      <a-form ref="formRef" :model="formState" layout="vertical" name="form_in_modal">
        <a-form-item
          name="social_name"
          label="社交名称"
          :rules="[{ required: true, message: '请输入社交名称' }]"
        >
          <a-input v-model:value="formState.social_name" />
        </a-form-item>
        <a-form-item
          name="social_value"
          label="社交值"
          :rules="[{ required: true, message: '请输入社交值' }]"
        >
          <a-input v-model:value="formState.social_value" />
        </a-form-item>
        <a-form-item name="css_class" label="图标样式">
          <a-radio-group v-model:value="formState.css_class">
            <a-radio :value="`i-fa6-brands:x-twitter`">
              <div class="i-fa6-brands:x-twitter w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:facebook`">
              <div class="i-fa6-brands:facebook w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:instagram`">
              <div class="i-fa6-brands:instagram w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:youtube`">
              <div class="i-fa6-brands:youtube w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:bilibili`">
              <div class="i-fa6-brands:bilibili w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:qq`">
              <div class="i-fa6-brands:qq w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:github`">
              <div class="i-fa6-brands:github w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:square-git`">
              <div class="i-fa6-brands:square-git w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:weixin`">
              <div class="i-fa6-brands:weixin w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-fa6-brands:zhihu`">
              <div class="i-fa6-brands:zhihu w-5 h-5"></div>
            </a-radio>
            <a-radio :value="`i-bi:link-45deg`">
              <div class="i-bi:link-45deg w-5 h-5"></div>
            </a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item
          label="是否为外链"
          name="is_link"
          class="collection-create-form_last-form-item"
        >
          <a-radio-group v-model:value="formState.is_link">
            <a-radio :value="false">否</a-radio>
            <a-radio :value="true">是</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
  <div>
    <a-table :columns="columns" :data-source="data">
      <template #bodyCell="{ column, text, record }">
        <template v-if="['social_name', 'social_value'].includes(column.dataIndex)">
          <div>
            <a-input
              v-if="editableData[record.id]"
              v-model:value="editableData[record.id][column.dataIndex]"
              style="margin: -5px 0"
            />
            <template v-else>
              {{ text }}
            </template>
          </div>
        </template>
        <template v-if="'css_class' === column.dataIndex">
          <div>
            <a-radio-group
              v-if="editableData[record.id]"
              v-model:value="editableData[record.id][column.dataIndex]"
            >
              <a-radio :value="`i-fa6-brands:x-twitter`">
                <div class="i-fa6-brands:x-twitter w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:facebook`">
                <div class="i-fa6-brands:facebook w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:instagram`">
                <div class="i-fa6-brands:instagram w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:youtube`">
                <div class="i-fa6-brands:youtube w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:bilibili`">
                <div class="i-fa6-brands:bilibili w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:qq`">
                <div class="i-fa6-brands:qq w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:github`">
                <div class="i-fa6-brands:github w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:square-git`">
                <div class="i-fa6-brands:square-git w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:weixin`">
                <div class="i-fa6-brands:weixin w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-fa6-brands:zhihu`">
                <div class="i-fa6-brands:zhihu w-5 h-5"></div>
              </a-radio>
              <a-radio :value="`i-bi:link-45deg`">
                <div class="i-bi:link-45deg w-5 h-5"></div>
              </a-radio>
            </a-radio-group>
            <template v-else>
              <div class="w-5 h-5" :class="getIcon(text)"></div>
            </template>
          </div>
        </template>
        <template v-if="column.key === 'is_link'">
          <a-radio-group
            v-if="editableData[record.id]"
            v-model:value="editableData[record.id][column.dataIndex]"
          >
            <a-radio :value="false">否</a-radio>
            <a-radio :value="true">是</a-radio>
          </a-radio-group>
          <template v-else>
            <a-tag color="success">{{ record.is_link ? '是' : '否' }}</a-tag>
          </template>
        </template>

        <template v-else-if="column.dataIndex === 'operation'">
          <div class="editable-row-operations">
            <span v-if="editableData[record.id]">
              <a-typography-link @click="save(record.id)">保存</a-typography-link>
              <a-popconfirm title="确定取消？" @confirm="cancel(record.id)">
                <a>取消</a>
              </a-popconfirm>
            </span>
            <span v-else>
              <a @click="edit(record.id)">编辑</a>
            </span>

            <a-popconfirm v-if="data.length" title="确认删除？" @confirm="deleteInfo(record.id)">
              <a>删除</a>
            </a-popconfirm>
          </div>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, type UnwrapRef } from 'vue'
import {
  AddSocial,
  DeleteSocial,
  GetSocial,
  type SocialConfig,
  type SocialConfigReq,
  UpdateSocial
} from '@/interfaces/Config'
import { type FormInstance, message } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'

const columns = [
  {
    title: '社交名称',
    dataIndex: 'social_name',
    key: 'social_name'
  },
  {
    title: '社交值',
    dataIndex: 'social_value',
    key: 'social_value'
  },
  {
    title: '图标样式',
    key: 'css_class',
    dataIndex: 'css_class'
  },
  {
    title: '是否为外链',
    key: 'is_link',
    dataIndex: 'is_link'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const data = ref<SocialConfig[]>([])

const getSocialInfo = async () => {
  try {
    const response: any = await GetSocial()
    data.value = response.data.data?.list || []
  } catch (error) {
    console.log(error)
  }
}
getSocialInfo()

// 添加社交信息
const formRef = ref<FormInstance>()
const visible = ref(false)
const formState = reactive<SocialConfigReq>({
  social_name: '',
  social_value: '',
  css_class: '',
  is_link: false
})

const addSocialInfo = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async (values) => {
        try {
          const response: any = await AddSocial(formState)
          if (response.data.code !== 0) {
            message.error(response.data.message)
            return
          }
          message.success('添加成功')
          visible.value = false
          if (formRef.value) {
            formRef.value.resetFields()
          }
          await getSocialInfo()
        } catch (error) {
          console.log(error)
        }
      })
      .catch((info) => {
        console.log('Validate Failed:', info)
        message.warning('请检查表单是否填写正确')
      })
  }
}

// 删除
const deleteInfo = async (id: string) => {
  try {
    const response: any = await DeleteSocial(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getSocialInfo()
  } catch (error) {
    console.log(error)
  }
}

// 编辑
const editableData: UnwrapRef<Record<string, SocialConfig>> = reactive({})
const edit = (id: string) => {
  editableData[id] = cloneDeep(data.value.filter((item) => id === item.id)[0])
}

const save = async (id: string) => {
  const editableDatum = editableData[id]
  try {
    const response: any = await UpdateSocial(id, editableDatum)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    await getSocialInfo()
  } catch (error) {
    console.log(error)
  }
}
const cancel = (key: string) => {
  delete editableData[key]
}

const getIcon = (icon: string): string => {
  switch (icon) {
    case 'i-fa6-brands:x-twitter':
      return 'i-fa6-brands:x-twitter'
    case 'i-fa6-brands:facebook':
      return 'i-fa6-brands:facebook'
    case 'i-fa6-brands:instagram':
      return 'i-fa6-brands:instagram'
    case 'i-fa6-brands:youtube':
      return 'i-fa6-brands:youtube'
    case 'i-fa6-brands:bilibili':
      return 'i-fa6-brands:bilibili'
    case 'i-fa6-brands:qq':
      return 'i-fa6-brands:qq'
    case 'i-fa6-brands:github':
      return 'i-fa6-brands:github'
    case 'i-fa6-brands:square-git':
      return 'i-fa6-brands:square-git'
    case 'i-fa6-brands:weixin':
      return 'i-fa6-brands:weixin'
    case 'i-fa6-brands:zhihu':
      return 'i-fa6-brands:zhihu'
    case 'i-bi:link-45deg':
      return 'i-bi:link-45deg'
  }
  return ''
}
</script>

<style scoped>
.collection-create-form_last-form-item {
  margin-bottom: 0;
}

.editable-row-operations a {
  margin-right: 8px;
}
</style>
