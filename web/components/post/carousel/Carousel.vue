<template>
  <div class="w-full h-85 relative b-rounded-4 overflow-hidden bg-#fff dark:text-dtc dark_bg_gray" ref="carouselRef"
       @mouseenter="addWheelListener" @mouseleave="removeWheelListener">
    <div v-if="carousel.length > 0" class="slides w-full h-full absolute flex transition-transform duration-700"
         :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
      <div class="slide w-full h-full relative flex-shrink-0" v-for="(item, index) in props.carousel"
           :key="index">
        <a
            class="relative h-full slide-up item group flex  b-rounded-4 cursor-pointer ease-linear duration-100 mb-5"
            :href="baseUrl + '/posts/' +item.id"
            target="_blank"
            :title="item.title"
        >
          <img
              class="w-full h-full"
              :src="apiHost + item.cover_img"
              :alt="item.title"
          />
          <div class="w-90% flex flex-col flex-center absolute  top-50% left-50% translate--50% translate--50%">
            <div class="text-6 font-bold">
              {{
                $dayjs(item.created_at * 1000).format("YYYY-MM-DD")
              }}
            </div>
            <div class="text-8 font-bold">
              {{ item.title }}
            </div>
            <div class="text-6 font-bold">
              {{ item.summary }}
            </div>
          </div>
        </a>
      </div>
    </div>
    <div v-else class="relative w-full h-full ">
      <div class="absolute top-50% left-50% translate--50% translate--50% text-10">
        轮播图暂无数据
      </div>
    </div>
    <!-- 圆点导航 -->
    <div class="absolute bottom-5 left-1/2 transform -translate-x-1/2 flex space-x-2">
      <button
          class="dot h-3 w-3 rounded-full bg-white"
          v-for="(slide, index) in carousel"
          :key="'dot-' + index"
          :class="{ 'bg-gray-800 w-7': currentSlide === index }"
          @click="setCurrentSlide(index)"
      ></button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {PropType} from "vue";
import type {CarouselVO} from "~/api/config";

const props = defineProps({
  carousel: {
    type: Array as PropType<CarouselVO[]>,
    default: () => [],
    required: true,
  },
});
const runtimeConfig = useRuntimeConfig()
const baseUrl = runtimeConfig.public.domain;
const apiHost = runtimeConfig.public.apiHost;

const currentSlide = ref(0);
const prevSlide = () => {
  currentSlide.value = (currentSlide.value - 1 + props.carousel.length) % props.carousel.length;
};
const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % props.carousel.length; // 确保索引在有效范围内循环
};

let slideInterval: number | undefined;
const intervalTime = 4000; // 间隔时间，例如 3000 毫秒（3秒）
const throttleInterval = 800; // 节流间隔时间
const startAutoSlide = () => {
  stopAutoSlide(); // 在开始新的定时器前确保清除已有的定时器
  slideInterval = setInterval(nextSlide, intervalTime) as unknown as number;
};

const stopAutoSlide = () => {
  if (slideInterval !== undefined) {
    clearInterval(slideInterval);
    slideInterval = undefined;
  }
};

onMounted(() => {
  startAutoSlide(); // 组件挂载后启动自动轮播
});

onUnmounted(() => {
  stopAutoSlide(); // 组件卸载前清除定时器
});

const setCurrentSlide = (index: number) => {
  currentSlide.value = index;
  startAutoSlide(); // 用户手动切换后重启自动轮播
};

let lastInvokeTime = 0; // 上次调用处理函数的时间
const throttle = (func: Function, delay: number) => {
  return (event: WheelEvent) => {
    event.preventDefault();
    const currentTime = Date.now();
    if (currentTime - lastInvokeTime > delay) {
      func(event);
      lastInvokeTime = currentTime;
    }
  };
};
const handleWheel = (event: WheelEvent) => {
  if (event.deltaY > 0) {
    nextSlide();
  } else if (event.deltaY < 0) {
    prevSlide();
  }
  startAutoSlide(); // 重启自动轮播，以保持用户体验一致性
};
const throttledHandleWheel = throttle(handleWheel, throttleInterval);

const carouselRef = ref<HTMLDivElement>();
const addWheelListener = () => {
  carouselRef.value?.addEventListener('wheel', throttledHandleWheel);
};

const removeWheelListener = () => {
  carouselRef.value?.removeEventListener('wheel', throttledHandleWheel);
};
</script>

<style scoped>
.dot {
  transition: background-color 0.3s;
}
</style>
