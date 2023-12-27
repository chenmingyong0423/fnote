<template>
  <div>
    <div class="flex items-center bg-white border-rounded-2">
      <input type="text" placeholder="input search text"  v-model="keyword" class="border-rounded-l-2 w-full w-90% h-[30px] bg-none light_border p-2 p-l-3 outline-#3F8DF1FF outline-1 border-r-none ">
      <div class=" w-10% light_border h-[30px] p-2 flex hover:border-#3F8DF1FF easy duration-100" @click="search">
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
})

const keyword = ref(props.keyword)

const search = () => {
  emit('search', keyword)
}

</script>