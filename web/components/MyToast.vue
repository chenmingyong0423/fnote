<template>
  <div
    v-if="toast.visible"
    class="fixed top-10 left-50% bg-#f44336 text-white p-y-5 p-x-4 b-rounded-4 drop-shadow-xl z-99999 slide-bottom-animation"
    :class="{ 'slide-up-animation': !toast.visible }"
  >
    {{ toast.message }}
  </div>
</template>

<script setup lang="ts">
import { useAlertStore } from "~/store/toast";

const toast = useAlertStore();
const isVisible = ref(false);

onUnmounted(() => {
  isVisible.value = false;
});
</script>

<style scoped>
@keyframes slideBottom {
  from {
    transform: translateY(-100%);
  }
  to {
    transform: translateY(0);
  }
}

.slide-bottom-animation {
  animation: slideBottom 0.3s ease-out;
}
</style>
