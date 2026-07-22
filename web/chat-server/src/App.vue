<template>
  <router-view />
</template>

<script>
import { onMounted } from "vue";
import { useAppStore } from "@/store";
import axios from "axios";
export default {
  name: "App",
  setup() {
    const store = useAppStore();
    const getUserInfo = async () => {
      try {
        const req = {
          uuid: store.userInfo.uuid,
        };
        const rsp = await axios.post(
          store.backendUrl + "/user/getUserInfo",
          req
        );
        if (rsp.data.code == 200) {
          if (!rsp.data.data.avatar.startsWith("http")) {
            rsp.data.data.avatar = store.backendUrl + rsp.data.data.avatar;
          }
          store.setUserInfo(rsp.data.data);
		  return rsp.data.data;
        } else {
          console.error(rsp.data.message);
		  return null;
        }
        console.log(rsp);
      } catch (error) {
        console.log(error);
		return null;
      }
    };
    onMounted(async () => {
      if (store.userInfo.uuid && store.token) {
        const currentUser = await getUserInfo();
        if (currentUser) {
          store.connectSocket();
        }
      }
    });
  },
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box; /* 推荐使用，以确保布局计算的一致性 */
}
</style>
