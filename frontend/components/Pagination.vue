<template>
    <!-- <div class="mt30 w200 mx-auto flex justify-between"> -->
    <div class="mt20">
        <el-pagination layout="prev,pager,next" :total="20" />
        <!-- <div
            class="w40 h40 bg-#1890ff/90 hover:bg-#1890ff/60 active:bg-#1890ff rounded-8 flex items-center justify-center cursor-pointer c-#fff">
            <div class="i-grommet-icons:previous text-22" />
        </div>
        <div
            class="w40 h40 bg-#1890ff/90 hover:bg-#1890ff/60 active:bg-#1890ff rounded-8 flex items-center justify-center cursor-pointer c-#fff">
            <div class="i-grommet-icons:next text-22" />
        </div> -->
    </div>
</template>


<script lang="ts" setup>

</script>

<style scoped>
:deep(.el-pagination) {
    justify-content: center;
    font-size: 24px !important;
}
</style>