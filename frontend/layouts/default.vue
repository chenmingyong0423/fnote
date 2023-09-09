<template>
    <div>
        <MyHeader class="fixed top-0 w-full z-99 " />
        <div class="bg-#F0F2F5 dark_bg_black pt100">
            <div class="max-w-1280 mx-auto  rounded-t-23 bg-#fff p30 min-h-100vh">
                <slot></slot>
            </div>
        </div>
        <MyFooter />
        <el-backtop :right="100" :bottom="100" />
    </div>
</template>
<script lang="ts" setup>

</script>

