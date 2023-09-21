<template>
    <div v-if="replies" class="mt25 px50">
        <el-space direction="vertical" :fill="true" class="w-full" :spacer="spacer">

            <div v-for="item in replies" :key="replies.id" class="b-1 b-#DDDBDB b-solid p20 rounded-10">
                <el-space :size="8" alignment="flex-start">
                    <el-avatar :size="36" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                    <div>
                        <el-space :size="0">
                            <div class="ml8 flex">
                                <div>{{ item.username }}</div>
                                <div v-if="item.replyTo" class="ml8">
                                    <el-space>
                                        <div>回复</div>
                                        <div>{{ item.replyTo }}</div>
                                    </el-space>
                                </div>
                            </div>
                            <div class="ml8">
                                发表于
                            </div>
                            <div class="ml4">
                                {{ item.replyTime }}
                            </div>
                        </el-space>
                    </div>
                </el-space>
                <v-md-preview :text="item.content" class="px18"></v-md-preview>
                <el-button class="ml50" @click="showReply = true">回复</el-button>
                <el-button class="ml50" v-if="showReply" @click="showReply = false">取消</el-button>
            </div>
        </el-space>
        <SubmitComments v-if="showReply" />
    </div>
</template>


<script lang="ts" setup>
import { ElDivider } from 'element-plus'
const props = defineProps(['replies'])
const spacer = h(ElDivider)
let showReply = ref(false)
</script>

<style scoped></style>