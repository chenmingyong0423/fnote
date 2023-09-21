<template>
    <div
        class="bg-#fff rounded-10 mt30 pb30 dark_bg_gray dark_text_white cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--5">
        <div class="text-16 pt20 px20">
            最新评论
        </div>
        <el-divider class="important:mt24" />
        <div class="text-14 px20 truncate" v-for="item in 5">
            <el-space direction="vertical" alignment="start" class="truncate">
                <el-text truncated class="dark_text_white ">
                    威：懂了。就是如果没有 cancel 信号，程序会在 select 处阻
                </el-text>
                <el-space class="c-#000/40 dark_text_white">
                    <div class="i-grommet-icons:article"></div>
                    <div>
                        一文掌握 Go 并发模式 Context 上下文
                    </div>
                </el-space>
            </el-space>
            <el-divider class="important:my10" />
        </div>
    </div>
</template>


<script lang="ts" setup>

</script>

<style scoped></style>