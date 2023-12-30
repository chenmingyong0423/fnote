<template>
  <div class="custom_border_gray text-center b-rounded-2 cursor-pointer select-none">
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