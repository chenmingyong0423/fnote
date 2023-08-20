<template>
    <div class="text-center text-30 py100 dark_text_white ">
        <div>这是列表页面</div>
        <div>名称:{{ route.params.name }}</div>
        <el-button @click="router.go(-1)">返回</el-button>
    </div>
</template>


<script lang="ts" setup>
const router = useRouter()
const route = useRoute()
</script>

<style scoped></style>