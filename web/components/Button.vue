<template>
  <div class="w-15 h-8 custom_border_gray text-center line-height-8 b-rounded-2 cursor-pointer">
    {{ props.name }}
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  name: {
    type: String,
    required: true,
  },
})
</script>