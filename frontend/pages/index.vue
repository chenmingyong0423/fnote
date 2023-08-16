<template>
    <div class="bg-#F0F2F5 py40 px25">
        <el-row :gutter="20">
            <el-col :span="17">
                <div>
                    <Content></Content>
                    <Content></Content>
                    <Content></Content>
                    <Content></Content>
                    <Content></Content>
                </div>
            </el-col>
            <el-col :span="7">
                <div class="w-full  ">
                    <Profile />
                    <Comment />
                </div>
            </el-col>
        </el-row>
    </div>
</template>


<script lang="ts" setup>

</script>

<style scoped></style>