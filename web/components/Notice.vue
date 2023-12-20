<template>
  <div class="flex items-center h-15 bg-#fff mb-5 b-rounded-4 cursor-pointer p-x-5">
    <div class="-i-ph-speaker-high-duotone w-10 h-10 text-orange"></div>
    <div ref="marqueeContainer" class="ml-5 font-bold overflow-hidden whitespace-nowrap w-full"
         @mouseenter="stopMarquee"
         @mouseleave="startMarquee"><span
        ref="marqueeContent"
        class="inline-block">{{ homeStore.notice_info.content }}</span>
      <span
          ref="marqueeContent2"
          class="inline-block m-l-5" v-show="showTheSecondMarquee">{{ homeStore.notice_info.content }}</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, onMounted} from 'vue';
import {useHomeStore} from "~/store/home";

const homeStore = useHomeStore()
const marqueeContainer = ref<HTMLElement | null>(null);
const marqueeContent = ref<HTMLElement | null>(null);
const marqueeContent2 = ref<HTMLElement | null>(null);
const showTheSecondMarquee = ref<boolean>(false);

const checkMarquee = () => {
  if (marqueeContent.value && marqueeContainer.value) {
    if (marqueeContent.value.offsetWidth > marqueeContainer.value.offsetWidth) {
      marqueeContent.value.classList.add('marquee-animation');
      showTheSecondMarquee.value = true
      marqueeContent2.value?.classList.add('marquee-animation');
    }
  }
};

const stopMarquee = () => {
  if (marqueeContent.value) {
    marqueeContent.value.style.animationPlayState = 'paused';
    if (marqueeContent2.value) {
      marqueeContent2.value.style.animationPlayState = 'paused';
    }
  }
};

const startMarquee = () => {
  if (marqueeContent.value) {
    marqueeContent.value.style.animationPlayState = 'running';
    if (marqueeContent2.value) {
      marqueeContent2.value.style.animationPlayState = 'running';
    }
  }
};

onMounted(checkMarquee);
</script>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(100%);
  }
  100% {
    transform: translateX(-100%);
  }
}

.marquee-animation {
  animation: marquee 20s linear infinite;
}
</style>