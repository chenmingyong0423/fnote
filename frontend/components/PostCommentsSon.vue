<template>
    <div v-if="replies" class="mt25 px50 lt-lg:px0">
        <el-space direction="vertical" :fill="true" class="w-full" :spacer="spacer">
            <div class="b-1 b-#DDDBDB b-solid p20 rounded-10">

                <!-- 回复信息 -->
                <el-space :size="8" alignment="flex-start">
                    <el-avatar :size="36" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                    <div class="text-14 lt-lg:text-14">
                        <!-- <el-space :size="0"> -->
                        <span class="ml8">
                            <span>{{ replies.name }}</span>
                            <span v-if="replies.reply_to" class="">
                                <span class="mx8">回复</span>
                                <span>{{ replies.reply_to }}</span>
                            </span>
                        </span>
                        <div class="display-none lt-lg:display-block lt-lg:my2" />
                        <span class="ml8">
                            发表于
                        </span>
                        <span class="ml4">
                            {{ formatTimestamp(replies.reply_time) }}
                        </span>
                        <!-- </el-space> -->
                    </div>
                </el-space>
                <!-- 回复内容 -->
                <v-md-preview :text="replies.content" class="px18 mt20"></v-md-preview>

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