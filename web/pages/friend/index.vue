<template>
  <div>
    <div class="bg-white b-rounded-4 p-5 mb-5 dark:text-dtc dark_bg_gray">
      <div class="h-10 line-height-10 font-bold text-6 mb-5">友链</div>
      <div class="flex flex-wrap gap-x-4 mb-5 lt-md:flex-col">
        <a
            :href="friend.url"
            target="_blank"
            class="flex mb-5 w-23% custom_border_gray h-[100px] p-2 b-rounded-4 custom_cursor_flow cursor-pointer dark:text-dtc dark_bg_gray custom_shadow lt-md:w-100%"
            v-for="(friend, index) in friends"
            :key="index"
        >
          <div class="w-15% flex justify-center">
            <img
                :src="friend.logo"
                alt=""
                class="w-15 h-15 border-rounded-50% cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
                v-if="friend.logo != ''"
            />
            <div
                class="i-ph-user-circle-duotone w-full h-12 border-rounded-50% lt-lg:mr0 text-gray-4"
                v-else
            ></div>
          </div>
          <div class="w-74% ml-1%">
            <div class="font-bold mb-2 mt-2">
              {{ friend.name }}
            </div>
            <div class="truncate">
              {{ friend.description }}
            </div>
          </div>
        </a>
      </div>
    </div>

    <div class="bg-white b-rounded-4 p-5 dark:text-dtc dark_bg_gray">
      <div class="text-10 mb-5">留言互友</div>
      <div class="flex flex-wrap gap-x-5 gap-y-2 lt-md:flex-col">
        <div class="flex w-49% lt-md:w-100%">
          <span
              class="light_border box-border w-[50px] border-rounded-l-2 custom_border_gray bg-#F9F9F9 h-10 line-height-10 text-center dark_bg_gray"
          >*昵称</span
          >
          <input
              type="text"
              placeholder="请输入昵称"
              v-model="req.name"
              class="w-full outline-none custom_border_gray box-border h-10 border-rounded-l-0 border-l-0 bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
          />
        </div>
        <div class="flex w-49% lt-md:w-100%">
          <span
              class="light_border box-border w-[50px] border-rounded-l-2 custom_border_gray bg-#F9F9F9 h-10 line-height-10 text-center dark_bg_gray"
          >邮箱</span
          >
          <input
              type="text"
              placeholder="请输入邮箱"
              v-model="req.email"
              class="w-full outline-none custom_border_gray box-border h-10 border-rounded-l-0 border-l-0 bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
          />
        </div>
        <div class="flex w-49% lt-md:w-100%">
          <span
              class="light_border box-border w-[100px] border-rounded-l-2 custom_border_gray bg-#F9F9F9 h-10 line-height-10 text-center dark_bg_gray"
          >*头像链接</span
          >
          <input
              type="text"
              placeholder="请输入头像链接"
              v-model="req.logo"
              class="w-full outline-none custom_border_gray box-border h-10 border-rounded-l-0 border-l-0 bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
          />
        </div>
        <div class="flex w-49% lt-md:w-100%">
          <span
              class="light_border box-border w-[100px] border-rounded-l-2 custom_border_gray bg-#F9F9F9 h-10 line-height-10 text-center dark_bg_gray"
          >*网站链接</span
          >
          <input
              type="text"
              placeholder="请输入网站链接"
              v-model="req.url"
              class="w-full outline-none custom_border_gray box-border h-10 border-rounded-l-0 border-l-0 bg-#F9F9F9 focus:custom_border_1E80FF b-rounded-2 p-2 box-border dark:text-dtc dark_bg_gray"
          />
        </div>
        <div class="w-100%">
          <textarea
              rows="5"
              class="w-full custom_border_gray bg-#F9F9F9 outline-none focus:custom_border_1E80FF b-rounded-2 p-2 box-border mb-3 dark:text-dtc dark_bg_gray"
              v-model="req.description"
              placeholder="*请输入个人简介（不能超出 30 字）"
              maxlength="200"
          ></textarea>
        </div>
      </div>
      <div>
        <Button
            name="提交"
            class="w-15 h-8 line-height-8 m-auto bg-#1E80FF text-white hover:bg-#1E80FF/70 duration-200"
            @click="submit"
        ></Button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import {
  type IFriend,
  type FriendReq,
  getFriends,
  applyForFriend,
} from "~/api/friend";
import {useAlertStore} from "~/store/toast";
import type {IBaseResponse, IListData, IResponse} from "~/api/http";
import {isValidEmail} from "~/utils/email";
import {useConfigStore} from "~/store/config";

const configStore = useConfigStore();
const friends = ref<IFriend[]>([]);
const toast = useAlertStore();

const req = ref<FriendReq>({
  name: "",
  email: "",
  logo: "",
  url: "",
  description: "",
});
const getFriendList = async () => {
  try {
    let httpRes: any = await getFriends();
    if (httpRes.data.value === null) {
      toast.showToast(httpRes.error.value.statusMessage, 2000);
      return;
    }
    let res: IResponse<IListData<IFriend>> = httpRes.data.value;
    if (res && res.data) {
      if (res.code !== 0) {
        toast.showToast(res.message, 2000);
        return;
      }
      friends.value = res.data?.list || [];
    }
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
};
await getFriendList();

const submit = async () => {
  try {
    if (req.value.name === "") {
      toast.showToast("请输入昵称", 2000);
      return;
    }
    if (req.value.logo === "") {
      toast.showToast("请输入头像链接", 2000);
      return;
    }
    if (req.value.url === "") {
      toast.showToast("请输入网站链接", 2000);
      return;
    }
    if (req.value.description === "") {
      toast.showToast("请输入个人简介", 2000);
      return;
    }
    if (req.value.email !== "") {
      if (!isValidEmail(req.value.email || "")) {
        toast.showToast("邮箱格式不正确！", 1000);
        return;
      }
    }
    const deepCopyReq: FriendReq = JSON.parse(JSON.stringify(req.value));
    let httpRes: any = await applyForFriend(deepCopyReq);
    if (httpRes.data.value === null) {
      if (httpRes.error.value.statusCode == 403) {
        toast.showToast("友链模块暂未开放！", 2000);
      } else if (httpRes.error.value.statusCode == 429) {
        toast.showToast("请勿重复提交！", 2000);
      } else {
        toast.showToast(httpRes.error.value.statusMessage, 2000);
      }
      return;
    }
    let res: IBaseResponse = httpRes.data.value;
    if (res) {
      if (res.code !== 0) {
        toast.showToast(res.message, 2000);
        return;
      }
      toast.showToast("提交成功，待站长审核通过后将会通过邮件告知。", 3000);
      req.value = {
        name: "",
        email: "",
        logo: "",
        url: "",
        description: "",
      };
    }
  } catch (error: any) {
    toast.showToast(error.toString(), 2000);
  }
};

useHead({
  title: `友链 - ${configStore.seo_meta_config.title === '' ? configStore.website_info.website_name : configStore.seo_meta_config.title}`,
  meta: [
    {name: "description", content: "友链列表"},
    {name: "keywords", content: configStore.seo_meta_config.keywords},
    {name: "author", content: configStore.seo_meta_config.author},
    {name: "robots", content: configStore.seo_meta_config.robots},
  ],
  link: [
    {rel: "icon", type: "image/x-icon", href: configStore.website_info.website_icon},
  ],
});
useSeoMeta({
  ogTitle: `友链 - ${configStore.seo_meta_config.og_title === '' ? configStore.website_info.website_name : configStore.seo_meta_config.og_title}`,
  ogDescription: "友链列表",
  ogImage: configStore.seo_meta_config.og_image,
  twitterCard: "summary",
});
</script>
