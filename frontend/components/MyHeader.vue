<template>
    <div class="flex justify-between items-center py5 dark_bg_gray" :class="homeStore.myHeaderBg">
        <div>
            <el-space>
                <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" :size="50"
                    class="mx30" />
                <el-menu router :default-active="activeIndex" mode="horizontal" @select="handleSelect" :ellipsis="false"
                    background-color="transparent" :text-color="homeStore.headerTextColor">
                    <el-menu-item index="/">首页</el-menu-item>
                    <el-menu-item index="/category">文章列表</el-menu-item>
                    <template v-for="item, index in homeStore.menuList.data.list" :key="index">
                        <!-- <el-sub-menu :index="String(index)" v-if="item.tags"
                            popper-class=" bg-#000/50 dark_bg_gray important:b-0">
                            <template #title>{{ item.name }}</template>
                            <el-menu-item v-for="tag, index in item.tags" :index="tag" :key="index">
                                {{ tag }}
                            </el-menu-item>
                        </el-sub-menu> -->
                        <el-menu-item :index="item.route">
                            {{ item.name }}
                        </el-menu-item>
                    </template>
                    <el-menu-item index="/friends">友链</el-menu-item>
                    <el-menu-item index="/about">关于我</el-menu-item>
                </el-menu>
            </el-space>
        </div>
        <div>
            <el-space :size="15">
                <el-switch size="large" v-model="isBlackMode" :active-action-icon="Moon" :inactive-action-icon="Sunny"
                    inactive-color="rgba(0,0,0,0.2)" active-color="#000" />
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
const isBlackMode = ref<boolean>(false)
const activeIndex = ref<string>('/')

watch(isBlackMode, (newValue) => {
    if (newValue)
        document.querySelector('html')!.classList.add("dark");
    else
        document.querySelector('html')!.classList.remove("dark");
})
const handleSelect = (key: string, keyPath: string[]) => {
    // console.log(key, keyPath)
}
const homeStore = useHomeStore()

onMounted(() => {
    window.addEventListener('scroll', () => {
        if (document.documentElement.scrollTop > 500) {
            homeStore.myHeaderBg = 'bg-#fff/80 backdrop-blur-60'
            homeStore.headerTextColor = '#000'
        }
        else {
            homeStore.myHeaderBg = 'bg-#000/20 backdrop-blur-20'
            homeStore.headerTextColor = '#fff'
        }
    })
})


</script>