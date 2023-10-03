<template>
    <div v-if="replies" class="mt25 px50">
        <el-space direction="vertical" :fill="true" class="w-full" :spacer="spacer">
            <div v-for="item in replies" :key="replies.id" class="b-1 b-#DDDBDB b-solid p20 rounded-10">
                <el-space :size="8" alignment="flex-start">
                    <el-avatar :size="36" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                    <div>
                        <el-space :size="0">
                            <div class="ml8 flex">
                                <div>{{ item.name }}</div>
                                <div v-if="item.reply_to" class="ml8">
                                    <el-space>
                                        <div>回复</div>
                                        <div>{{ item.reply_to }}</div>
                                    </el-space>
                                </div>
                            </div>
                            <div class="ml8">
                                发表于
                            </div>
                            <div class="ml4">
                                {{ formatTimestamp(item.reply_time) }}
                            </div>
                        </el-space>
                    </div>
                </el-space>
                <div>
                    <v-md-preview :text="item.content" class="px18"></v-md-preview>
                <div v-if=" item.replied_content "
                class="p-10 whitespace-nowrap overflow-hidden overflow-ellipsis bg-gray-200 rounded-lg px18">
                    {{ item.replied_content }}
                </div>
                </div>
                <el-button class="ml50" @click="showReply = true">回复</el-button>
                <el-button class="ml50" v-if="showReply" @click="showReply = false">取消</el-button>
            </div>
        </el-space>
        <SubmitComments v-if="showReply" />
    </div>
</template>


<script lang="ts" setup>
import { ElDivider } from 'element-plus'
import { IReply } from '~/api/comment'
import { dayjs } from 'element-plus'
import { useHomeStore } from '../store/home';
const info = useHomeStore()

const props = defineProps(['replies'])
const replies: IReply[] = props.replies
const spacer = h(ElDivider)
let showReply = ref(false)
const formatTimestamp = (timestamp: number): string => {
    return dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss');
};

if (replies != undefined && replies.length > 0) {
    const myMap: Map<string, IReply> = new Map();
    for (let i = 0; i < replies.length; i++) {
        if (replies[i].name === info.masterInfo.name) {
            let temp = replies[i]
            temp.name = temp.name + '[作者]'
            replies[i] = temp
        }
        myMap.set(replies[i].id, replies[i])
    }
    for (let i = 0; i < replies.length; i++) {
        if (replies[i].reply_to_id !== "") {
            let rl = myMap.get(replies[i].reply_to_id)
            if (rl !== undefined) {
                let temp = replies[i]
                temp.replied_content = rl.content
                replies[i] = temp
            }
        }
    }
}
</script>

<style scoped></style>