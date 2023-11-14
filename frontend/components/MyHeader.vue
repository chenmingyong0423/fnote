<template>
    <div ref="myHeader"
        class="bg-#fff/20  backdrop-blur-20 fixed top-0 w-full z-99 flex justify-between items-center  dark_bg_gray duration-200 ease-linear">

        <div>
            <el-space>
                <!-- 菜单栏头像 -->
                <el-avatar :src="picture" :size="36"
                    class="mx30 cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
                    @click="router.push('/')" />
                <!-- 菜单项 -->
                <div class="flex">
                    <div class="menu_item" @click="router.push('/')" :class="route.path === '/' ? 'active' : ''">
                        首页
                    </div>
                    <div class="menu_item lt-lg:display-none " @click="router.push('/category')"
                        :class="route.path === '/category' ? 'active' : ''">
                        文章列表
                    </div>
                    <div class="menu_item lt-lg:display-none" @click="router.push(item.route)"
                        v-for="item in homeStore.menuList" :class="route.path === item.route ? 'active' : ''">
                        {{ item.name }}
                    </div>
                    <div class="menu_item lt-lg:display-none" @click="router.push('/friends')"
                        :class="route.path === '/friends' ? 'active' : ''">
                        友链
                    </div>
                    <div class="menu_item lt-lg:display-none" @click="router.push('/about')"
                        :class="route.path === '/about' ? 'active' : ''">
                        关于我
                    </div>
                    <!-- 移动端下拉菜单 -->
                    <div class="menu_item display-none lt-lg:display-block">
                        <el-dropdown trigger="click" size="large">
                            <div class="i-grommet-icons:more text-24 c-#000" />
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item @click="router.push('/category')">
                                        文章列表
                                    </el-dropdown-item>
                                    <el-dropdown-item v-for="item in homeStore.menuList" @click="router.push(item.route)">
                                        {{ item.name }}
                                    </el-dropdown-item>
                                    <el-dropdown-item @click="router.push('/friends')">
                                        友链
                                    </el-dropdown-item>
                                    <el-dropdown-item @click="router.push('/about')">
                                        关于我
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </div>
                </div>
            </el-space>
        </div>
        <div class="">
            <el-space :size="15">
                <el-switch size="large" v-model="homeStore.isBlackMode" :active-action-icon="Moon"
                    :inactive-action-icon="Sunny" inactive-color="rgba(0,0,0,0.2)" active-color="#000"
                    lt-lg:important:display-none />
                <div @click="homeStore.searchVisible = true"
                    class=" bg-#1890ff rounded-50% py6 px6 dark_bg_black cursor-pointer hover:bg-#4EA6F9 active:bg-#C7C7C8">
                    <div class="i-grommet-icons:search text-24 c-#fff " />
                </div>
                <div
                    class="i-grommet-icons:github text-36 dark_text_white cursor-pointer hover:bg-#999 active:bg-#e5e5e5 lt-lg:display-none" />
            </el-space>
        </div>
        <SearchDialog />
    </div>
</template>

<script lang="ts" setup>
import { useHomeStore } from '~/store/home';
import { Sunny, Moon } from '@element-plus/icons-vue'
const homeStore = useHomeStore()
const isBlackMode = computed(() => homeStore.isBlackMode)
const router = useRouter()
const route = useRoute()


watch(isBlackMode, (newValue) => {
    if (newValue)
        document.querySelector('html')!.classList.add("dark");
    else
        document.querySelector('html')!.classList.remove("dark");
})
const handleSelect = (key: string, keyPath: string[]) => {
    console.log(key, keyPath)
    // router.push(key)
}
const myHeader = ref()
let scrollCount = ref(0)
const headerScroll = () => {

    // const scrollNow = ref(document.documentElement.scrollTop)
    // console.log(Math.abs(scrollNow.value - scrollCount.value));
    if (document.documentElement.scrollTop > scrollCount.value) {
        myHeader.value.setAttribute('style', `top:-${myHeader.value.clientHeight}px`);
    }
    else {
        myHeader.value.setAttribute('style', 'top:0px');
    }
    scrollCount.value = document.documentElement.scrollTop
}
onMounted(() => {
    window.addEventListener('scroll', headerScroll)
})
onBeforeUnmount(() => {
    window.removeEventListener('scroll', headerScroll)
})

const picture = homeStore.masterInfo.picture


import { getMenus, IMenu } from "../api/category"
import { IResponse, IListData } from "../api/http";
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

</script>
<style scoped>
.active {
    background-color: rgba(0, 0, 0, 0.2);
}

:global(.el-dropdown-menu--large .el-dropdown-menu__item) {
    font-size: 16px;
}
</style>