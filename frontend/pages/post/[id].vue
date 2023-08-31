<template>
    <div class="text-center text-30 py100 dark_text_white ">
        <div>这是文章页面</div>
        <div>文章id:{{ route.params.id }}</div>
        <el-button @click="router.go(-1)">返回</el-button>
    </div>
</template>


<script lang="ts" setup>
const router = useRouter()
const route = useRoute()
const data = {
    "code": 200,
    "message": "OK",
    "data": {
        "sug": "post1",
        "author": "陈明勇",
        "title": "Go语言啊",
        "summary": "Summary 1",
        "cover_img": "/images/cover1.jpg",
        "category": "A",
        "tags": [
            "a",
            "b"
        ],
        "likeCount": 2,
        "comments": 3,
        "visit": 0,
        "priority": 1,
        "createTime": 1692806408145,
        "meta_description": "Description 1",
        "meta_keywords": "Keyword 1",
        "word_count": 0,
        "allow_comment": true,
        "update_time": 1692806408143
    }
}
</script>

<style scoped></style>