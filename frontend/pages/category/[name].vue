<template>
    <div>
        <div class="py15">
            <div class="text-30 mb20 ml10">分类:{{ homeStore.classification ? homeStore.classification?.name : '未知' }}</div>
        </div>

        <div class="py15">
            <div class="text-30 mb20 ml10">
                标签:
                <el-space>
                    <my-tag :show-close-btn="true" @click="popTags(tag)" v-for="tag, index in activeTags" :key="index">
                        {{ tag }}
                    </my-tag>
                </el-space>
            </div>
            <el-space>
                <my-tag v-for="tag, index in tags" :key="index" @click="getPostByTag(tag)" mt20>
                    {{ tag }}
                </my-tag>
            </el-space>
        </div>

        <div class="bg-#e5e5e5/40 p20 rounded-10" v-if="dataList.length > 0">
            <div v-for="item, index in dataList" :key="index">
                <Content @click="router.push(`/post/${item.sug}`)" :postData="item"></Content>
                <!-- <el-divider /> -->
            </div>
        </div>
        <el-empty description="暂无数据" v-else />
        <Pagination />
    </div>
</template>

<script lang="ts" setup>
import { useHomeStore } from '~/store/home';
import { IPost, getPosts, PageRequest } from '~/api/post';
import { IResponse, IPageData } from "~/api/http";
import { IMenu } from "~/api/category"
import { getTagList } from '~/api/tag'
const route = useRoute()
const router = useRouter()
const homeStore = useHomeStore()
const tags = ref([])//标签列表
const activeTags = ref<string[]>([])//当前选中的标签列表
homeStore.classification = homeStore.menuList.find((item: IMenu) => {
    if (item.route === route.path)
        return item
})

console.log(homeStore.classification);

const dataList = ref<IPost[]>([])
const rq = ref<PageRequest>({
    pageNo: 1,
    pageSize: 5,
    sortField: '',
    sortOrder: '',
    search: '',
    category: '',
    tags: [`${homeStore.classification ? homeStore.classification!.name : null}`]
})

const postInfos = async () => {
    try {
        let postRes: any = await getPosts(rq.value)
        let res: IResponse<IPageData<IPost>> = postRes.data.value
        dataList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
}
const getTags = async () => {
    const result: any = await getTagList(homeStore.classification!.name)
    tags.value = result.data.value.data.list
}
getTags()
postInfos()

watch(activeTags.value, (newValue) => {
    console.log(newValue);
    rq.value.tags = activeTags.value
    postInfos()
})
// 数据穿梭
const sendData = (data: string, originList: string[], TargetData: string[]) => {
    const index = originList.findIndex((item) => data === item)
    console.log(index);
    TargetData.push(data)
    originList.splice(index, 1)
}
// 点击标签查找
const getPostByTag = (name: string) => {
    sendData(name, tags.value, activeTags.value)
}
//删除选中标签
const popTags = (name: string) => {
    sendData(name, activeTags.value, tags.value)
}
</script>

<style scoped></style>
