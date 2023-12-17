<template>
    <div
        class="dark_text_white w-auto inline-block  px16 py8 mx10 text-18 b-2 b-#1890ff b-solid rounded-12 cursor-pointer hover:b-#4A87EF hover:bg-#e5e5e5/30 active-bg-#e5e5e5 active-c-#fff">

        <div>
            <span>#</span>
            <slot></slot>
            <div class="i-grommet-icons:close  c-#bbb text-16 ml10 inline-block" v-if="showCloseBtn"
                @click="closeCallback" />
        </div>


    </div>
</template>


<script lang="ts" setup>
defineProps({
    'closeCallback': {
        type: (): any => { },
        default: () => { }
    },
    'showCloseBtn': {
        type: Boolean,
        default: false
    }
})
</script>

<style scoped></style>