import { defineStore } from "pinia";

export const useAlertStore = defineStore("toast", {
  state: () => ({
    message: "",
    visible: false,
  }),
  actions: {
    showToast(msg: string, duration: number) {
      this.message = msg;
      this.visible = true;
      setTimeout(() => this.hideToast(), duration);
    },
    hideToast() {
      this.visible = false;
    },
  },
});
