<template>
  <div v-if="visible" class="fixed z-999 inset-0 bg-black bg-opacity-40 flex items-center justify-center p-4">
    <div class="bg-white p-6 rounded shadow-lg max-w-md min-w-400px dark:text-dtc dark_bg_full_black">
      <span @click="closeModal" class="float-right text-xl cursor-pointer">&times;</span>
      <p class="leading-loose text-5">{{ content }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">

const emit = defineEmits(['tooltipVisibleChanged']);

const props = defineProps({
  content: {
    type: String,
    required: true,
  },
  visible: {
    type: Boolean,
    required: true,
  },
});

const closeModal = () => {
  emit('tooltipVisibleChanged', false);
};

</script>

<style>
</style>
