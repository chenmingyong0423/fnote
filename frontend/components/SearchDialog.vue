<template>
    <el-dialog :append-to-body="true" v-model="homeStore.searchVisible" :lock-scroll="false" title="搜索文章" width="85%">
        <div class="mx20">
            <div class="flex">
                <el-input v-model="md" placeholder="请输入搜索内容" />
                <div
                    class=" bg-#1890ff w30 h30 flex justify-center items-center dark_bg_black cursor-pointer hover:bg-#4EA6F9 active:bg-#C7C7C8">
                    <div class="i-grommet-icons:search text-14 c-#fff " />
                </div>
            </div>
            <el-empty description="暂无数据" />
        </div>
    </el-dialog>
</template>


<script lang="ts" setup>
import { useHomeStore } from '~/store/home';
const homeStore = useHomeStore()
const md = ref<any>()//文本框输入搜索内容
</script>

<style scoped></style>