<template>
    <!-- bg-#F0F2F5 -->
    <div>
        <!-- 后期加 -->
        <!-- <div class="text-0 relative">
            <div class="text-center w-full font-bold text-70 absolute c-#fff top-30% text-shadow-xl">
                xxx's Blog
                <TypeIt></TypeIt>
            </div>
            <img src="../assets/images/bg.png" class="w-full h-100vh object-cover" />
        </div> -->

        <div class="pt70 pb40 px25">
            <el-row :gutter="20">
                <el-col :span="17">
                    <div>
                        <Content v-for="item in data.data.list " :postData="item">
                        </Content>
                    </div>
                </el-col>
                <el-col :span="7">
                    <div class="w-full">
                        <Profile />
                        <el-affix :offset="0">
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
const data = {
    "code": 200,
    "message": "OK",
    "data": {
        "list": [
            {
                "sug": "post5",
                "author": "陈明勇",
                "title": "哈哈",
                "summary": "Summary 1",
                "cover_img": "/images/cover1.jpg",
                "category": "D",
                "tags": [
                    "tag1",
                    "tag2"
                ],
                "likeCount": 2,
                "comments": 3,
                "visit": 0,
                "priority": 1,
                "createTime": 1692806408149
            },
            {
                "sug": "post5",
                "author": "陈明勇",
                "title": "哈哈",
                "summary": "Summary 1",
                "cover_img": "/images/cover1.jpg",
                "category": "D",
                "tags": [
                    "tag1",
                    "tag2"
                ],
                "likeCount": 2,
                "comments": 3,
                "visit": 0,
                "priority": 1,
                "createTime": 1692806408149
            },
            {
                "sug": "post5",
                "author": "陈明勇",
                "title": "哈哈",
                "summary": "Summary 1",
                "cover_img": "/images/cover1.jpg",
                "category": "D",
                "tags": [
                    "tag1",
                    "tag2"
                ],
                "likeCount": 2,
                "comments": 3,
                "visit": 0,
                "priority": 1,
                "createTime": 1692806408149
            },
            {
                "sug": "post5",
                "author": "陈明勇",
                "title": "哈哈",
                "summary": "Summary 1",
                "cover_img": "/images/cover1.jpg",
                "category": "D",
                "tags": [
                    "tag1",
                    "tag2"
                ],
                "likeCount": 2,
                "comments": 3,
                "visit": 0,
                "priority": 1,
                "createTime": 1692806408149
            },
            {
                "sug": "post5",
                "author": "陈明勇",
                "title": "哈哈",
                "summary": "Summary 1",
                "cover_img": "/images/cover1.jpg",
                "category": "D",
                "tags": [
                    "tag1",
                    "tag2"
                ],
                "likeCount": 2,
                "comments": 3,
                "visit": 0,
                "priority": 1,
                "createTime": 1692806408149
            }
        ]
    }
}
onMounted(() => {


})
definePageMeta({
    layout: "home"
})
</script>

<style scoped></style>