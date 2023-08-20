<template>
    <!-- bg-#F0F2F5 -->
    <div>
        <div class="text-0 relative">
            <div class="text-shadow-xl text-center w-full font-bold text-70 absolute c-#fff top-30% ">
                xxx's Blog
            </div>
            <img src="../assets/images/bg.png" class="w-full h-100vh object-cover" />
        </div>

        <div class=" py40 px25">
            <el-row :gutter="20">
                <el-col :span="17">
                    <div>
                        <Content @click="router.push(`detail/${5}`)" v-for="item in 5 "></Content>
                    </div>
                </el-col>
                <el-col :span="7">
                    <div class="w-full">
                        <Profile />
                        <el-affix :offset="50">
                            <Comment />
                        </el-affix>
                    </div>
                </el-col>
            </el-row>
        </div>
    </div>
</template>


<script lang="ts" setup>
const router = useRouter()
onMounted(() => {

})
</script>

<style scoped></style>