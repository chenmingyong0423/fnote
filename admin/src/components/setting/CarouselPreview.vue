<template>
  <div class="carousel-preview">
    <img v-if="coverImg" :src="imageUrl" :alt="title" />
    <div class="carousel-overlay" :style="overlayStyle">
      <div class="carousel-badge">精选文章</div>
      <div class="carousel-title" :style="textStyle" :title="title">
        {{ title || '轮播图标题' }}
      </div>
      <div class="carousel-summary" :style="textStyle">
        {{ summary || '轮播图摘要' }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = defineProps<{
  coverImg: string
  title: string
  summary: string
  color: string
  serverHost: string
}>()

const resolvedColor = computed(() =>
  /^#[0-9a-f]{6}$/i.test(props.color) ? props.color : '#ffffff'
)

const isDarkText = computed(() => {
  const color = resolvedColor.value
  const red = Number.parseInt(color.slice(1, 3), 16)
  const green = Number.parseInt(color.slice(3, 5), 16)
  const blue = Number.parseInt(color.slice(5, 7), 16)
  return (red * 0.299 + green * 0.587 + blue * 0.114) / 255 < 0.48
})

const imageUrl = computed(() =>
  /^https?:\/\//i.test(props.coverImg) ? props.coverImg : props.serverHost + props.coverImg
)

const overlayStyle = computed(() => ({
  background: isDarkText.value
    ? 'linear-gradient(to top, rgba(255, 255, 255, 0.85), rgba(255, 255, 255, 0.5), transparent)'
    : 'linear-gradient(to top, rgba(0, 0, 0, 0.75), rgba(0, 0, 0, 0.4), transparent)'
}))

const textStyle = computed(() => ({
  color: resolvedColor.value,
  textShadow: isDarkText.value
    ? '0 1px 3px rgba(255, 255, 255, 0.9)'
    : '0 1px 3px rgba(0, 0, 0, 0.8)'
}))
</script>

<style scoped>
.carousel-preview {
  position: relative;
  width: 100%;
  height: 256px;
  overflow: hidden;
  background: #f5f5f5;
  border-radius: 8px;
}

.carousel-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-overlay {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  padding: 48px 20px 24px;
}

.carousel-badge {
  display: inline-flex;
  margin-bottom: 8px;
  padding: 4px 10px;
  color: #fff;
  font-size: 11px;
  font-weight: 500;
  line-height: 16px;
  background: rgba(0, 0, 0, 0.45);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 999px;
  backdrop-filter: blur(4px);
}

.carousel-title {
  overflow: hidden;
  font-size: 20px;
  font-weight: 700;
  line-height: 1.375;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.carousel-summary {
  display: -webkit-box;
  margin-top: 6px;
  overflow: hidden;
  font-size: 14px;
  line-height: 20px;
  opacity: 0.9;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

@media (max-width: 576px) {
  .carousel-preview {
    height: 176px;
  }

  .carousel-overlay {
    padding: 48px 16px 20px;
  }

  .carousel-title {
    font-size: 16px;
  }

  .carousel-summary {
    font-size: 12px;
    line-height: 16px;
  }
}
</style>
