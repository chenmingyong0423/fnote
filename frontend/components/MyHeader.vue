<template>
    <div class="flex justify-between items-center bg-#000/20 backdrop-blur-20 py5 dark_bg_gray ">
        <div>
            <el-space>
                <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" :size="50"
                    class="mx30" />
                <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" @select="handleSelect"
                    :ellipsis="false" background-color="transparent" text-color="#fff">
                    <el-menu-item index="index">首页</el-menu-item>
                    <el-menu-item index="list">大列表页</el-menu-item>
                    <template v-for="item, index in menuList.data.list" :key="index">
                        <el-menu-item :index="String(index)" v-if="!item?.tags">
                            {{ item.name }}
                        </el-menu-item>
                        <el-sub-menu :index="String(index)" v-else popper-class=" bg-#000/50 dark_bg_gray important:b-0">
                            <template #title>{{ item.name }}</template>
                            <el-menu-item v-for="tag, index in item?.tags" :index="tag" :key="index">
                                {{ tag }}
                            </el-menu-item>
                        </el-sub-menu>
                    </template>
                    <el-menu-item index="friends">友链</el-menu-item>
                    <el-menu-item index="about">关于我</el-menu-item>
                </el-menu>
            </el-space>
        </div>
        <div>
            <el-space :size="15">
                <el-switch size="large" v-model="value1" :active-action-icon="Moon" :inactive-action-icon="Sunny"
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
import { Sunny, Moon } from '@element-plus/icons-vue'
const value1 = ref(false)
const activeIndex = ref('index')
const menuList = {
    "code": 200,
    "message": "OK",
    "data": {
        "list": [
            {
                "name": "后端",
                "route": "/category/backend"
            },
            {
                "name": "前端",
                "route": "/category/frontend"
            }
        ]
    }
}
watch(value1, () => {
    if (value1.value == true) {
        document.querySelector('html')!.classList.add("dark");
    }
    else {
        document.querySelector('html')!.classList.remove("dark");
    }
})
const handleSelect = (key: string, keyPath: string[]) => {
    console.log(key, keyPath)
}
</script>