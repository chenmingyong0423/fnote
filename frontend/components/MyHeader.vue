<template>
    <div class="flex justify-between items-center px30 py5">
        <div>
            <el-space>
                <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" :size="50"
                    class="mx30" />
                <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" @select="handleSelect"
                    :ellipsis="false">
                    <el-menu-item index="1">首页</el-menu-item>
                    <el-sub-menu index="2">
                        <template #title>技术文章</template>
                        <el-menu-item index="2-1">item one</el-menu-item>
                        <el-menu-item index="2-2">item two</el-menu-item>
                        <el-menu-item index="2-3">item three</el-menu-item>
                    </el-sub-menu>
                    <el-menu-item index="3">文章专栏</el-menu-item>
                    <el-menu-item index="4">开源库</el-menu-item>
                    <el-menu-item index="5">友链</el-menu-item>
                    <el-menu-item index="6">关于我</el-menu-item>
                </el-menu>
            </el-space>
        </div>
        <div>
            <el-space :size="15">
                <div class=" bg-blue rounded-50% py6 px6">
                    <div class="i-grommet-icons:search text-24 c-#fff " />
                </div>

                <div class="i-grommet-icons:github text-36" />
            </el-space>
        </div>
    </div>
</template>

<script lang="ts" setup>
const activeIndex = ref('1')
const handleSelect = (key: string, keyPath: string[]) => {
    console.log(key, keyPath)
}
</script>