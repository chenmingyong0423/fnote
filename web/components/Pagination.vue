<template>
  <div class="flex bg-white items-center gap-x-4 p-2 b-rounded-2 w-full dark:text-dtc dark_bg_gray">
    <div class="list-none flex gap-x-5 m-auto">
      <NuxtLink class="w-[32px] h-[32px] flex-center"
                :class="{'text-gray-3 select-none' : currentPage == 1, 'cursor-pointer hover:bg-gray-1': currentPage != 1}"
                :to="route+ calculatePages(currentPage - 1) + '?pageSize=' + perPageSize"
      >
        &lt
      </NuxtLink>
      <NuxtLink class="w-[32px] h-[32px] cursor-pointer hover:bg-gray-1 flex-center"
                v-if="maxPageNumbers - currentPage <= 1"
                :to="route+ 1 + '?pageSize=' + perPageSize"
      >
        1
      </NuxtLink>
      <NuxtLink class="w-[32px] h-[32px] text-gray-3 flex-center select-none cursor-pointer hover:text-#1e80ff"
                v-if="maxPageNumbers - currentPage <= 0"
                @mouseover="leftPointer = '<<'"
                @mouseleave="leftPointer = '...'"
                :to="route + calculatePages(currentPage - maxPageNumbers) + '?pageSize=' + perPageSize"
      >
        {{ leftPointer }}
      </NuxtLink>

      <NuxtLink class="w-[32px] h-[32px] cursor-pointer hover:bg-gray-1 flex-center"
                v-for="(page, index) in pagesToShow"
                :class="{ 'border-solid border-#1e80ff border-1 b-rounded-2': page === currentPage }"
                :key="index"
                :to="route + page + '?pageSize=' + perPageSize"
      >
        {{ page }}
      </NuxtLink>
      <NuxtLink
          class="w-[32px] h-[32px] text-gray-3 flex-center select-none hover:text-#1e80ff cursor-pointer hover:text-#1e80ff"
          v-if="totalPages > maxPageNumbers && !pagesToShow.includes(totalPages)"
          @mouseover="rightPointer = '>>'"
          @mouseleave="rightPointer = '...'"
          :to="route + calculatePages(currentPage + maxPageNumbers) + '?pageSize=' + perPageSize"
      >
        {{ rightPointer }}
      </NuxtLink>
      <NuxtLink class="w-[32px] h-[32px] cursor-pointer hover:bg-gray-1 flex-center"
                v-if="!pagesToShow.includes(totalPages)"
                :to="route+ totalPages + '?pageSize=' + perPageSize"
      >{{ totalPages }}
      </NuxtLink>
      <NuxtLink class="w-[32px] h-[32px] flex-center"
                :class="{'text-gray-3 select-none' : currentPage == totalPages, 'cursor-pointer hover:bg-gray-1': currentPage != totalPages}"
                :to="route+ calculatePages(currentPage + 1) + '?pageSize=' + perPageSize"
      >
        &gt
      </NuxtLink>
      <div>
        <select class="w-[95px] h-[35px] p-x-2 p-y-1 b-rounded-2 border-gray-3 dark:text-dtc dark_bg_gray " v-model="perPageSize" @change="changeItemsPerPage">
          <option class="dark:bg-black" :value="1">1 条/页</option>
          <option class="dark:bg-black" :value="5">5 条/页</option>
          <option class="dark:bg-black" :value="10">10 条/页</option>
          <option class="dark:bg-black" :value="20">20 条/页</option>
          <option class="dark:bg-black" :value="50">50 条/页</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>

const leftPointer = ref("...");
const rightPointer = ref("...");

const props = defineProps({
  total: {
    type: Number,
    required: true
  },
  currentPage: {
    type: Number,
    default: 2
  },
  perPageCount: {
    type: Number,
    default: 10
  },
  route: {
    type: String,
    required: true
  },
});

const route = props.route
const currentPage = ref(props.currentPage);
const perPageSize = ref(props.perPageCount);
// 最大显示多少个页码
const maxPageNumbers = ref(5);
const total = ref(props.total)
const totalPages = computed(() => Math.ceil(total.value / perPageSize.value));

const pagesToShow = computed(() => {
  let pages = [];
  let start = Math.max(currentPage.value - Math.floor(maxPageNumbers.value / 2), 1);
  let end = Math.min(start + maxPageNumbers.value - 1, totalPages.value);

  if (totalPages.value > maxPageNumbers.value && end === totalPages.value) {
    start = end - maxPageNumbers.value + 1;
  }

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  return pages;
});

watch(() => props.total, () => {
  if (props.total != total.value) {
    total.value = props.total
  }
});

const router = useRouter()

const changeItemsPerPage = () => {
  router.push(route + 1 + '?pageSize=' + perPageSize.value)
};

const calculatePages = (page: number): number => {
  if (page < 1) {
    return 1;
  } else if (page > totalPages.value) {
    return totalPages.value;
  }
  return page
};
</script>