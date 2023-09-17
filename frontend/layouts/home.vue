<template>
    <div>
        <MyHeader />
        <div class="bg-#F0F2F5 dark_bg_black min-h-100vh">
            <slot></slot>
        </div>
        <MyFooter />
        <el-backtop :right="100" :bottom="100" />
    </div>
</template>
<script lang="ts" setup>

</script>

