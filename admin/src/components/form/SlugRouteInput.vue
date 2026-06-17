<template>
  <a-input
    :value="value"
    :placeholder="placeholder"
    @input="handleInput"
    @update:value="emitValue"
  />
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { normalizeSlug, shouldAutoFillSlug } from '@/utils/slug'

const props = withDefaults(
  defineProps<{
    value: string
    source: string
    placeholder?: string
  }>(),
  {
    placeholder: '例如 tech-note'
  }
)

const emit = defineEmits<{
  'update:value': [value: string]
}>()

const touched = ref(false)

const handleInput = () => {
  touched.value = true
}

const emitValue = (value: string) => {
  emit('update:value', value)
}

const resetTouched = () => {
  touched.value = false
}

watch(
  () => props.source,
  (source) => {
    if (!touched.value) {
      emit('update:value', shouldAutoFillSlug(source) ? normalizeSlug(source) : '')
    }
  }
)

defineExpose({
  resetTouched
})
</script>
