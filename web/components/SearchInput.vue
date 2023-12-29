<template>
  <div>
    <div class="flex items-center bg-white border-rounded-2 dark:text-dtc dark_bg_gray">
      <input type="text" placeholder="input search text"  v-model="keyword" class="outline-none border-rounded-l-2 custom_border_gray w-full w-90% h-[30px] p-2 p-l-3 dark:text-dtc dark_bg_gray">
      <div class=" w-10% light_border h-[30px] p-2 flex hover:border-#3F8DF1FF easy duration-100 dark:text-dtc dark_bg_gray" @click="search">
        <span class="i-ph:magnifying-glass text-gray w-full h-full cursor-pointer"></span>
      </div>
      <span></span>
    </div>
  </div>
</template>

<script lang="ts" setup>

const emit = defineEmits(['search']);

const props = defineProps({
  keyword: {
    type: String,
    default: "",
    required: false,
  },
  callbackRoutePath: {
    type: String,
    default: "/search",
    required: false,
  }
})

const keyword = ref(props.keyword)

const router = useRouter()

const search = () => {
  router.push(props.callbackRoutePath + '?keyword=' + keyword.value)
}

</script>