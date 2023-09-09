<template>
    <div>
        <div class="py15">
            <div class="text-30 mb20 ml10">标签:{{ tag[0]?.name || '未知' }}</div>
        </div>
        <div>
            <div v-for="item in data.data.list ">
                <Content @click="router.push(`/post/${item.sug}`)" :postData="item"></Content>
                <el-divider />
            </div>
        </div>
        <div class="pl500">
            <el-pagination layout="prev, pager, next" :total="5000" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useHomeStore } from '~/store/home';

const route = useRoute()
const router = useRouter()
const homeStore = useHomeStore()

const tag = homeStore.menuList.data.list.filter((item) => {
    if (item.route === route.path)
        return item
})

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
            }
        ]
    }
}
</script>

<style scoped></style>