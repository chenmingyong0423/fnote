<template>
  <div
      class="bg-#fff  backdrop-blur-20 fixed top-0 dark_text_gray w-full z-99 flex justify-between items-center p-1 dark_bg_gray duration-200 ease-linear">
    <div>
      <a-avatar :size="50" :src="picture"
                class="mx5 cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"></a-avatar>
    </div>
    <div class="bg-transparent">
      <a-menu v-model:selectedKeys="current" mode="horizontal" @click="menuItemChanged" :items="items"
              class="dark_text_gray bg-transparent"/>
    </div>
    <div class="ml-auto pr-5">
      <a-space>
        暗黑模式
        <a-switch v-model:checked="homeStore.isBlackMode" checked-children="开" un-checked-children="关"/>
        <div
            class=" bg-#1890ff rounded-50% py2 px2 dark_bg_black cursor-pointer hover:bg-#4EA6F9 active:bg-#C7C7C8">
          <div class="i-grommet-icons:search text-5 c-#fff "/>
        </div>
        <div
            class="i-grommet-icons:github text-10 dark_text_white cursor-pointer hover:bg-#999 active:bg-#e5e5e5 lt-lg:display-none"/>
      </a-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {h, ref, onMounted} from 'vue';
import {FileMarkdownOutlined, HomeOutlined, TeamOutlined, UserOutlined, FileWordOutlined} from "@ant-design/icons-vue";
import type {MenuProps, ItemType} from "ant-design-vue";
import {getMenus} from "~/server/api/category"
import type {IMenu} from "~/server/api/category"
import type {IResponse, IListData} from "~/server/api/http";
import {useHomeStore} from '~/store/home';

const route = useRoute()
const router = useRouter()

const current = ref<string[]>([route.path])
const items = ref<MenuProps['items']>([]);
const homeStore = useHomeStore()

const isBlackMode = computed(() => homeStore.isBlackMode)
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

onMounted(() => {
  items.value?.push({
        key: '/index',
        icon: () => h(HomeOutlined),
        label: '首页',
        title: '首页',
      },
      {
        key: '/list',
        icon: () => h(FileMarkdownOutlined),
        label: '文章列表',
        title: '文章列表',
      },
      {
        key: '/friend',
        icon: () => h(TeamOutlined),
        label: '友链',
        title: '友链',
      },
      {
        key: '/about-me',
        icon: () => h(UserOutlined),
        label: '关于我',
        title: '关于我',
      }
  )
  let i = ref<number>(2)
  homeStore.menuList.forEach((item: IMenu) => {
    items.value?.splice(i.value, 0, {
      key: item.route,
      icon: () => h(FileWordOutlined),
      label: item.name,
      title: item.name,
    })
    i.value++
  })
})

watch(isBlackMode, (newValue) => {
  localStorage.setItem("isBlackMode", newValue.toString())
  if (newValue)
    document.querySelector('html')!.classList.add("dark");
  else
    document.querySelector('html')!.classList.remove("dark");
})
</script>
