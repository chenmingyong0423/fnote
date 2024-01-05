<template>
  <div>
    <div
      @click="visible = true"
      class="flex items-center dark:text-dtc dark_bg_gray h-15 bg-#fff mb-5 b-rounded-4 cursor-pointer p-x-5 ease-linear duration-100 custom_shadow hover:translate-y--2"
    >
      <div class="-i-ph-speaker-high-duotone w-10 h-10 text-orange"></div>
      <div
        ref="marqueeContainer"
        class="ml-5 font-bold overflow-hidden whitespace-nowrap w-full"
        @mouseenter="stopMarquee"
        @mouseleave="startMarquee"
      >
        <span ref="marqueeContent" class="inline-block">
          [{{
            $dayjs(homeStore.notice_info.publish_time * 1000).format(
              "YYYY-MM-DD",
            )
          }}] {{ homeStore.notice_info.title }}</span
        >
        <span
          ref="marqueeContent2"
          class="inline-block m-l-5"
          v-show="showTheSecondMarquee"
          >{{ homeStore.notice_info.title }}</span
        >
      </div>
    </div>
    <!-- 模态框 -->
    <div
      v-if="visible"
      class="fixed z-999 inset-0 bg-black bg-opacity-40 flex items-center justify-center p-4 custom_shadow"
    >
      <div
        class="bg-white p-6 rounded shadow-lg max-w-md min-w-400px dark:text-dtc dark_bg_full_black"
      >
        <div class="text-right text-sm text-gray-500 mb-4">
          发布时间:
          {{
            $dayjs(homeStore.notice_info.publish_time * 1000).format(
              "YYYY-MM-DD HH:mm:ss",
            )
          }}
        </div>
        <span @click="closeModal" class="float-right text-xl cursor-pointer"
          >&times;</span
        >
        <h2 class="text-xl font-bold mb-4">
          {{ homeStore.notice_info.title }}
        </h2>
        <p class="indent-8 leading-loose">
          {{ homeStore.notice_info.content }}
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useHomeStore } from "~/store/home";

const homeStore = useHomeStore();
const marqueeContainer = ref<HTMLElement | null>(null);
const marqueeContent = ref<HTMLElement | null>(null);
const marqueeContent2 = ref<HTMLElement | null>(null);
const showTheSecondMarquee = ref<boolean>(false);
const visible = ref(false);

const closeModal = () => {
  visible.value = false;
};
const checkMarquee = () => {
  if (marqueeContent.value && marqueeContainer.value) {
    if (marqueeContent.value.offsetWidth > marqueeContainer.value.offsetWidth) {
      marqueeContent.value.classList.add("marquee-animation");
      showTheSecondMarquee.value = true;
      marqueeContent2.value?.classList.add("marquee-animation");
    }
  }
};

const stopMarquee = () => {
  if (marqueeContent.value) {
    marqueeContent.value.style.animationPlayState = "paused";
    if (marqueeContent2.value) {
      marqueeContent2.value.style.animationPlayState = "paused";
    }
  }
};

const startMarquee = () => {
  if (marqueeContent.value) {
    marqueeContent.value.style.animationPlayState = "running";
    if (marqueeContent2.value) {
      marqueeContent2.value.style.animationPlayState = "running";
    }
  }
};

onMounted(checkMarquee);
</script>

<style scoped></style>
