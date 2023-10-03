<template>
    <div>
        <div
            class="text-center bg-#fff rounded-10 py25 dark_bg_gray dark_text_white cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--5">
            <el-avatar :size="100" :src="data.picture"
                class="hover:rotate-360 ease-out duration-1000 hover:drop-shadow-xl" />
            <div class="text-14 c-#000/60 mt35 dark_text_white">
                <el-space direction="vertical">
                    <div>{{ data.name }}</div>
                    <div>{{ data.profile }}</div>
                </el-space>
            </div>
            <el-divider />
            <div class="text-14">
                <el-row>
                    <el-col :span="8">
                        <div>
                            <el-space direction="vertical" :size="5">
                                <div>{{ data.post_count }}</div>
                                <div>文章</div>
                            </el-space>
                        </div>
                    </el-col>
                    <el-col :span="8">
                        <div class="w-full b-x-1 b-y-0 b-#ccc/60 b-solid">
                            <el-space direction="vertical" :size="5">
                                <div>{{ data.category_count }}</div>
                                <div>文章分类</div>
                            </el-space>
                        </div>
                    </el-col>
                    <el-col :span="8">
                        <el-space direction="vertical" :size="5">
                            <div>{{ data.website_views }}</div>
                            <div>访问量</div>
                        </el-space>
                    </el-col>
                </el-row>
            </div>
        </div>
    </div>
</template>


<script lang="ts" setup>
import { useHomeStore } from '../store/home';
const globalData = useHomeStore()
const data = globalData.masterInfo
</script>

<style scoped></style>