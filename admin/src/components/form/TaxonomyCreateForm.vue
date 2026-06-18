<template>
  <a-form ref="formRef" :model="formModel" layout="vertical" name="taxonomy_create_form">
    <a-form-item
      name="name"
      :label="nameLabel"
      :rules="[{ required: true, message: `请输入${nameLabel}` }]"
    >
      <a-input :value="name" :placeholder="`请输入${nameLabel}`" @update:value="updateName" />
    </a-form-item>
    <a-form-item
      name="route"
      label="前端路由"
      :rules="taxonomyRouteRules"
      :extra="taxonomyRouteExtra"
    >
      <SlugRouteInput
        ref="routeInputRef"
        :value="route"
        :source="name"
        :placeholder="routePlaceholder"
        @update:value="updateRoute"
      />
    </a-form-item>
    <template v-if="type === 'category'">
      <a-form-item name="description" label="描述">
        <a-textarea
          :value="description"
          placeholder="可选"
          allow-clear
          @update:value="updateDescription"
        />
      </a-form-item>
      <a-form-item name="show_in_nav" label="显示在导航">
        <a-switch :checked="showInNav" @update:checked="updateShowInNav" />
      </a-form-item>
    </template>
    <a-form-item name="enabled" label="启用">
      <a-switch :checked="enabled" @update:checked="updateEnabled" />
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import SlugRouteInput from '@/components/form/SlugRouteInput.vue'
import { createSlugRule } from '@/utils/slug'

const props = withDefaults(
  defineProps<{
    type: 'category' | 'tag'
    name: string
    route: string
    enabled: boolean
    description?: string
    showInNav?: boolean
  }>(),
  {
    description: '',
    showInNav: false
  }
)

const emit = defineEmits<{
  'update:name': [value: string]
  'update:route': [value: string]
  'update:description': [value: string]
  'update:enabled': [value: boolean]
  'update:showInNav': [value: boolean]
}>()

const formRef = ref<FormInstance>()
const routeInputRef = ref<InstanceType<typeof SlugRouteInput>>()

const nameLabel = computed(() => (props.type === 'category' ? '分类名称' : '标签名称'))
const routePlaceholder = computed(() => (props.type === 'category' ? '例如 tech-note' : '例如 vue'))
const taxonomyRouteExtra = computed(
  () => `用于前台 URL，建议使用小写英文、数字和短横线，${routePlaceholder.value}`
)
const taxonomyRouteRules = [createSlugRule('前端路由')]
const formModel = computed(() => ({
  name: props.name,
  route: props.route,
  description: props.description,
  show_in_nav: props.showInNav,
  enabled: props.enabled
}))

const updateName = (value: string) => {
  emit('update:name', value)
}

const updateRoute = (value: string) => {
  emit('update:route', value)
}

const updateDescription = (value: string) => {
  emit('update:description', value)
}

const updateEnabled = (value: boolean) => {
  emit('update:enabled', value)
}

const updateShowInNav = (value: boolean) => {
  emit('update:showInNav', value)
}

const validateFields = () => {
  return formRef.value?.validateFields() ?? Promise.resolve()
}

const resetFields = () => {
  formRef.value?.resetFields()
  routeInputRef.value?.resetTouched()
}

const resetRouteTouched = () => {
  routeInputRef.value?.resetTouched()
}

const clearValidate = () => {
  formRef.value?.clearValidate()
}

defineExpose({
  validateFields,
  resetFields,
  resetRouteTouched,
  clearValidate
})
</script>
