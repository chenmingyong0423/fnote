<template>
  <div
    v-if="visible"
    class="fixed inset-0 bg-black bg-opacity-40 flex justify-center items-center z-9999"
    @click="close"
  >
    <div class="relative max-w-90vw max-h-90vh" @click.stop>
      <img :src="imageUrl" alt="Image Preview" class="max-w-full max-h-full" />
    </div>
    <button
      class="absolute top-2 right-2 text-3xl text-white bg-transparent border-none cursor-pointer"
      @click="close"
    >
      ×
    </button>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true
  },
  imageUrl: {
    type: String,
    required: true
  }
})

const visible = ref<Boolean>(props.modelValue)

const emit = defineEmits(['update:modelValue'])

watch(
  () => props.modelValue,
  (value: Boolean) => {
    visible.value = value
  }
)

const close = () => {
  visible.value = false
  emit('update:modelValue', false)
}
</script>

<style scoped>
/* 你可以在这里添加特定于这个组件的样式，如果有需要的话 */
</style>
