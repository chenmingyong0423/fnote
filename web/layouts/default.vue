<template>
  <div class="dark:bg-#03080c">
    <SmallMenu></SmallMenu>
    <MyToast></MyToast>
    <div ref="myDom">
      <Header class="slide-down" />
    </div>
    <div class="bg-#F0F2F5 dark:bg-#03080c pt-25 p-5">
      <div class="w-90% m-auto lt-md:w-99%">
        <slot></slot>
      </div>
    </div>
    <div>
      <Footer />
    </div>
    <div>
      <div
        ref="scrollToTop"
        class="i-ph:caret-circle-up-fill text-#1e80ff dark:text-dtc w-15 h-15 fixed bottom-[100px] right-[20px] cursor-pointer"
        style="display: none"
        @click="handleScrollToTop"
      ></div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { useHomeStore } from "~/store/home";
import {
  type IWebsiteInfo,
  getWebsiteInfo,
  type TPSVVO,
  type IWebsite,
  type INotice,
  type ISocialInfo,
  type IPayInfo,
  type SeoMetaConfigVO,
} from "~/api/config";
import { type IResponse } from "~/api/http";
import SmallMenu from "~/components/SmallMenu.vue";
import {
  collectVisitLog,
  getWebsiteCountStats,
  type VisitLogRequest,
  type WebsiteCountStats,
} from "~/api/statiscs";
import { useConfigStore } from "~/store/config";
import type { IPost } from "~/api/post";

const myDom = ref();
const homeStore = useHomeStore();
const configStore = useConfigStore();
const runtimeConfig = useRuntimeConfig();
onMounted(() => {
  let isBlackMode = localStorage.getItem("isBlackMode");
  if (isBlackMode === "") {
    localStorage.setItem("isBlackMode", homeStore.isBlackMode.toString());
  } else {
    homeStore.isBlackMode = isBlackMode === "true";
  }
});

const siteName = ref("fnote");
const siteURL = runtimeConfig.public.domain;
const apiHost = runtimeConfig.public.apiHost;
const metaVerificationList = ref([] as { name: string; content: string }[]);

const { data: webMaster } = await useAsyncData<IWebsiteInfo>(
  `webMaster`,
  async () => {
    try {
      let postRes: any = await getWebsiteInfo();
      let res: IResponse<IWebsiteInfo> = postRes.data.value;
      if (res && res.data) {
        return res.data;
      } else {
        return {
          website_config: {} as IWebsite,
          notice_config: {} as INotice,
          social_info_config: {} as ISocialInfo,
          pay_info_config: [] as IPayInfo[],
          seo_meta_config: {} as SeoMetaConfigVO,
          third_party_site_verification: [] as TPSVVO[],
        };
      }
    } catch (error) {
      console.log(error);
      return {
        website_config: {} as IWebsite,
        notice_config: {} as INotice,
        social_info_config: {} as ISocialInfo,
        pay_info_config: [] as IPayInfo[],
        seo_meta_config: {} as SeoMetaConfigVO,
        third_party_site_verification: [] as TPSVVO[],
      };
    }
  },
);

configStore.website_info = webMaster.value?.website_config || ({} as IWebsite);
siteName.value = webMaster.value?.website_config.website_name || "fnote";
configStore.notice_info = webMaster.value?.notice_config || ({} as INotice);
configStore.social_info_list =
  webMaster.value?.social_info_config.social_info_list || [];
configStore.pay_info = webMaster.value?.pay_info_config || ([] as IPayInfo[]);
configStore.seo_meta_config =
  webMaster.value?.seo_meta_config || ({} as SeoMetaConfigVO);
configStore.tpsv_list =
  webMaster.value?.third_party_site_verification || ([] as TPSVVO[]);
configStore.tpsv_list.forEach((item) => {
  metaVerificationList.value.push({
    name: item.key,
    content: item.value,
  });
});

const { data: websiteCountStats } = await useAsyncData<WebsiteCountStats>(
  `websiteCountStats`,
  async () => {
    try {
      let postRes: any = await getWebsiteCountStats();
      let res: IResponse<WebsiteCountStats> = postRes.data.value;
      if (res && res.data) {
        return res.data;
      } else {
        return {} as WebsiteCountStats;
      }
    } catch (error) {
      console.log(error);
      return {} as WebsiteCountStats;
    }
  },
);

configStore.website_count_stats =
  websiteCountStats.value || ({} as WebsiteCountStats);

const scrollToTop = ref();

const scrollEvent = () => {
  if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 20) {
    scrollToTop.value.style.display = "block";
  } else {
    scrollToTop.value.style.display = "none";
  }
};
const handleScrollToTop = () => {
  if (process.client) {
    // 客户端特有的代码
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  }
};

const jsonLd = computed(() => {
  return JSON.stringify({
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: `${siteName.value}`,
    url: `${siteURL}`,
  });
});

useHead({
  script: [
    {
      type: "application/ld+json",
      innerHTML: jsonLd.value,
    },
  ],
  title:
    configStore.seo_meta_config.title === ""
      ? configStore.website_info.website_name
      : configStore.seo_meta_config.title,
  meta: [
    {
      name: "description",
      content: configStore.seo_meta_config.description || "fnote",
    },
    {
      name: "keywords",
      content: configStore.seo_meta_config.keywords || "fnote",
    },
    { name: "author", content: configStore.seo_meta_config.author || "fnote" },
    { name: "robots", content: configStore.seo_meta_config.robots || "fnote" },
    ...metaVerificationList.value,
  ],
  link: [
    {
      rel: "icon",
      type: "image/x-icon",
      href: apiHost + configStore.website_info.website_icon,
    },
  ],
});
useSeoMeta({
  ogTitle:
    configStore.seo_meta_config.og_title === ""
      ? configStore.website_info.website_name
      : configStore.seo_meta_config.title,
  ogDescription: configStore.seo_meta_config.description || "fnote",
  ogImage: configStore.seo_meta_config.og_image,
  twitterCard: "summary",
});
onMounted(() => {
  window.addEventListener("scroll", scrollEvent);
});

onBeforeUnmount(() => {
  window.removeEventListener("scroll", scrollEvent);
});

const collect = async () => {
  if (process.client) {
    // 客户端特有的代码
    try {
      const req = {
        url: window.location.href,
      } as VisitLogRequest;
      await collectVisitLog(req);
    } catch (error) {
      console.log(error);
    }
  }
};
collect();
</script>

<style scoped>
@keyframes slideDown {
  0% {
    transform: translateY(-100%);
  }
  100% {
    transform: translateY(0);
  }
}

.slide-down {
  animation: slideDown 1s ease;
  animation-iteration-count: 1;
}
</style>
