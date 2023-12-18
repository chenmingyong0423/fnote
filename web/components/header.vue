<template>
  <div
      class="bg-#fff  backdrop-blur-20 fixed top-0 dark:text-dtc w-full z-99 flex justify-between items-center p-1 dark_bg_gray duration-200 ease-linear">
    <div>
      <a-avatar :size="50" :src="picture"
                class="mx5 cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"></a-avatar>
    </div>
    <div class="bg_transparent">
<!--      <a-menu v-model:selectedKeys="current"  mode="horizontal" @click="menuItemChanged" :items="items"-->
<!--              class="dark:text-dtc bg_transparent" />-->
      <Menu :items="homeStore.menuList"></Menu>
    </div>
    <div class="ml-auto pr-5">
      <a-space>
        <div class="i-ph-sun-dim-duotone dark:i-ph-moon-stars-fill cursor-pointer text-10 text-#86909c dark:text-dtc dark:hover:text-white" @click="homeStore.is_black_mode = !homeStore.is_black_mode"></div>
        <div class="i-ph-list-magnifying-glass-duotone cursor-pointer dark:text-dtc text-10 text-#86909c dark:hover:text-white"></div>

        <div
            class="i-grommet-icons:github text-12 dark:text-dtc cursor-pointer dark:hover:text-white active:bg-#e5e5e5 lt-lg:display-none"/>
      </a-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {h, ref, onMounted, reactive} from 'vue';
import type {MenuProps, ItemType} from "ant-design-vue";
import {getMenus} from "~/server/api/category"
import type {IMenu} from "~/server/api/category"
import type {IResponse, IListData} from "~/server/api/http";
import {useHomeStore} from '~/store/home';
import Menu from "~/components/menu.vue";

const router = useRouter()
const homeStore = useHomeStore()

const isBlackMode = computed(() => homeStore.is_black_mode)
const picture = ref<string>(homeStore.master_info.picture)
const menus = async () => {
  try {
    let postRes: any = await getMenus()
    let res: IResponse<IListData<IMenu>> = postRes.data.value
    homeStore.menuList = res.data?.list || []
  } catch (error) {
    console.log(error);
  }
};
menus()

const menuItemChanged = (item: ItemType) => {
  router.push(item?.key as string)
}

watch(isBlackMode, (newValue) => {
  localStorage.setItem("isBlackMode", newValue.toString())
  if (newValue)
    document.querySelector('html')!.classList.add("dark");
  else
    document.querySelector('html')!.classList.remove("dark");
})
</script>

<style scoped>
</style>
