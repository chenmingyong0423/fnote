<template>
  <div class="flex bg-white items-center gap-x-4 p-2 b-rounded-2 w-full">
    <ul class="list-none flex gap-x-5 m-auto">
      <li class="w-[32px] h-[32px] flex-center"
          :class="{'text-gray-3 select-none' : currentPage == 1, 'cursor-pointer hover:bg-gray-1': currentPage != 1}"
          @click="changePage(currentPage - 1)"
      >
        &lt
      </li>
      <li class="w-[32px] h-[32px] cursor-pointer hover:bg-gray-1 flex-center" v-if="maxPageNumbers - currentPage <= 1"
          @click="changePage(1)">1
      </li>
      <li class="w-[32px] h-[32px] text-gray-3 flex-center select-none cursor-pointer hover:text-#1e80ff"
          v-if="maxPageNumbers - currentPage <= 0"
          @mouseover="leftPointer = '<<'"
          @mouseleave="leftPointer = '...'"
          @click="changePage(currentPage - maxPageNumbers)"
      >
        {{ leftPointer }}
      </li>

      <li class="w-[32px] h-[32px] cursor-pointer hover:bg-gray-1 flex-center"
          v-for="(page, index) in pagesToShow"
          :class="{ 'border-solid border-#1e80ff border-1 b-rounded-2': page === currentPage }"
          @click="changePage(page)"
          :key="index"
      >
        {{ page }}
      </li>
      <li class="w-[32px] h-[32px] text-gray-3 flex-center select-none hover:text-#1e80ff cursor-pointer hover:text-#1e80ff"
          v-if="totalPages > maxPageNumbers && !pagesToShow.includes(totalPages)"
          @mouseover="rightPointer = '>>'"
          @mouseleave="rightPointer = '...'"
          @click="changePage(currentPage + maxPageNumbers)"

      >
        {{ rightPointer }}
      </li>
      <li class="w-[32px] h-[32px] cursor-pointer hover:bg-gray-1 flex-center" v-if="!pagesToShow.includes(totalPages)"
          @click="changePage(totalPages)">{{ totalPages }}
      </li>
      <li class="w-[32px] h-[32px] flex-center"
          :class="{'text-gray-3 select-none' : currentPage == totalPages, 'cursor-pointer hover:bg-gray-1': currentPage != totalPages}"
          @click="changePage(currentPage + 1)"
      >
        &gt
      </li>
      <li>
        <select class="w-[95px] h-[35px] p-x-2 p-y-1 b-rounded-2 border-gray-3" v-model="perPageSize"
                @change="changeItemsPerPage">
          <option :value="1">1 条/页</option>
          <option :value="5">5 条/页</option>
          <option :value="10">10 条/页</option>
          <option :value="20">20 条/页</option>
          <option :value="50">50 条/页</option>
        </select>
      </li>
    </ul>
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
    default: 1
  },
  perPageCount: {
    type: Number,
    default: 10
  },
});
const emit = defineEmits(['pageChanged', 'perPageChanged']);

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
  currentPage.value = 1;
  total.value = props.total
});

watch(pagesToShow, (newValue, oldValue) => {
  console.log(newValue);
});

const changePage = (page: number) => {
  if (page < 1) {
    page = 1;
  } else if (page > totalPages.value) {
    page = totalPages.value;
  }
  if (page != currentPage.value) {
    currentPage.value = page;
    emit('pageChanged', page, perPageSize.value);
  }
};

const changeItemsPerPage = () => {
  currentPage.value = 1
  emit('perPageChanged', perPageSize.value);
};
</script>