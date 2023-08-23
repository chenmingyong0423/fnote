<template>
    <div class="pt100 ">
        <div class="max-w-1280 mx-auto  rounded-t-23 bg-#fff p30">
            <el-breadcrumb :separator-icon="ArrowRight">
                <el-breadcrumb-item :to="{ path: '/' }">homepage</el-breadcrumb-item>
                <el-breadcrumb-item>promotion management</el-breadcrumb-item>
                <el-breadcrumb-item>promotion list</el-breadcrumb-item>
                <el-breadcrumb-item>promotion detail</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="mt30">
                <el-check-tag v-for="tag in tags" :key="tag.name" class="mx5" checked :type="tag.type">
                    {{ tag.name }}
                </el-check-tag>
            </div>
            <div>
                <Content @click="router.push(`post/${item}`)" v-for="item in 5 "></Content>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ArrowRight } from '@element-plus/icons-vue'
const tags = ref([
    { name: '前端', type: '' },
    { name: '后端', type: 'success' },
    { name: 'Java', type: 'info' },
    { name: 'JavaScript', type: 'warning' },
    { name: 'python', type: 'danger' },
])
const router = useRouter()
</script>

<style scoped></style>