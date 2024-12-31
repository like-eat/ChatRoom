<template>
  <router-view />
</template>

<script>
import { onMounted } from 'vue';
import { useStore } from "vuex";
export default {
  name: "App",
  setup() {
    const store = useStore();
    onMounted(() => {
      if (store.state.userInfo.uuid) {
        const wsUrl =
          store.state.wsUrl + "/ws?client_id=" + store.state.userInfo.uuid;
        store.state.socket = new WebSocket(wsUrl);
        store.state.socket.onopen = () => {
          console.log("WebSocket连接已打开");
        };
        store.state.socket.onmessage = (message) => {
          console.log("收到消息：", message.data);
        };
        store.state.socket.onclose = () => {
          console.log("WebSocket连接已关闭");
        };
        store.state.socket.onerror = () => {
          console.log("WebSocket连接发生错误");
        };
        console.log(store.state.socket);
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