<template>
    <div>
        <div class="py15">
            <div class="text-30 mb20 ml10">标签</div>
            <my-tag v-for="tag in homeStore.menuList.data.list" :key="tag.name" @click="router.push(tag.route)">
                {{ tag.name }}
            </my-tag>
        </div>
        <div class="mt40">
            <div class="text-30 mb20 ml10">我的文章</div>
            <div v-for="item in data.data.list ">
                <Content @click="router.push(`post/${item}`)" :postData="item"></Content>
                <el-divider />
            </div>

        </div>
        <div class="">
            <el-pagination layout="prev, pager, next" :total="5000" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useHomeStore } from '~/store/home';
const homeStore = useHomeStore()
const tags = ref([
    { name: '前端', type: '' },
    { name: '后端', type: 'success' },
    { name: 'Java', type: 'info' },
    { name: 'JavaScript', type: 'warning' },
    { name: 'python', type: 'danger' },
])
const router = useRouter()
const data = {
    "code": 200,
    "message": "OK",
    "data": {
        "PageNo": 1,
        "PageSize": 1,
        "totalPages": 5,
        "totalCount": 5,
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
            }
        ]
    }
}
</script>

<style scoped>
.el-pagination {
    justify-content: center;
}
</style>