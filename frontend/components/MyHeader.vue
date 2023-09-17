<template>
    <div ref="myHeader"
        class="fixed top-0 w-full z-99 flex justify-between items-center dark_bg_gray duration-200 ease-linear"
        :class="homeStore.myHeaderBg">
        <!-- <div>
            <el-space>
                <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" :size="50"
                    class="mx30" />
                <el-menu router :default-active="$route.path" mode="horizontal" @select="handleSelect" :ellipsis="false"
                    background-color="transparent" :text-color="homeStore.headerTextColor">
                    <el-menu-item index="/">首页</el-menu-item>
                    <el-menu-item index="/category">文章列表</el-menu-item>
                    <template v-for="item, index in homeStore.menuList.data.list" :key="index">
                        <el-menu-item :index="item.route">
                            {{ item.name }}
                        </el-menu-item>
                    </template>
                    <el-menu-item index="/friends">友链</el-menu-item>
                    <el-menu-item index="/about">关于我</el-menu-item>
                </el-menu>
            </el-space>
        </div> -->
        <div>
            <el-space>
                <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" :size="36"
                    class="mx30 cursor-pointer hover:rotate-360 ease-out duration-1000" @click="router.push('/')" />
                <div class="flex">
                    <div class="menu_item" @click="router.push('/')" :class="route.path === '/' ? 'active' : ''">
                        首页
                    </div>
                    <div class="menu_item" @click="router.push('/category')"
                        :class="route.path === '/category' ? 'active' : ''">
                        文章列表
                    </div>
                    <div class="menu_item" @click="router.push(item.route)" v-for="item in homeStore.menuList.data.list"
                        :class="route.path === item.route ? 'active' : ''">
                        {{ item.name }}
                    </div>
                    <div class="menu_item" @click="router.push('/friends')"
                        :class="route.path === '/friends' ? 'active' : ''">
                        友链
                    </div>
                    <div class="menu_item" @click="router.push('/about')" :class="route.path === '/about' ? 'active' : ''">
                        关于我
                    </div>
                </div>
            </el-space>
        </div>
        <div>
            <el-space :size="15">
                <el-switch size="large" v-model="homeStore.isBlackMode" :active-action-icon="Moon"
                    :inactive-action-icon="Sunny" inactive-color="rgba(0,0,0,0.2)" active-color="#000" />
                <div class=" bg-#1890ff rounded-50% py6 px6 dark_bg_black">
                    <div class="i-grommet-icons:search text-24 c-#fff " />
                </div>
                <div class="i-grommet-icons:github text-36 dark_text_white" />
            </el-space>
        </div>
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
        myHeader.value.setAttribute('style', 'top:-50px');
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

</script>
<style scoped>
.active {
    background-color: rgba(0, 0, 0, 0.2);

}
</style>