import { createPinia, defineStore } from 'pinia'

export const pinia = createPinia()

export const useAppStore = defineStore('app', {
  state: () => ({
    backendUrl: process.env.VUE_APP_BACKEND_URL || 'http://127.0.0.1:8000',
    wsUrl: process.env.VUE_APP_WS_URL || 'ws://127.0.0.1:8000',
    userInfo: (sessionStorage.getItem('userInfo') && JSON.parse(sessionStorage.getItem('userInfo'))) || {},
    token: sessionStorage.getItem('token') || '',
    socket: null,
  }),
  actions: {
    setUserInfo(userInfo) {
      this.userInfo = userInfo;
      sessionStorage.setItem('userInfo', JSON.stringify(userInfo));
    },
    setToken(token) {
      this.token = token;
      sessionStorage.setItem('token', token);
    },
    connectSocket() {
      if (!this.token) {
        return;
      }
      if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
        return;
      }
      const wsUrl = this.wsUrl + '/wss';
      this.socket = new WebSocket(wsUrl, ['kama-chat', this.token]);
      this.socket.onopen = () => console.log('WebSocket连接已打开');
      this.socket.onmessage = (message) => console.log('收到消息：', message.data);
      this.socket.onclose = () => console.log('WebSocket连接已关闭');
      this.socket.onerror = (error) => console.error('WebSocket连接发生错误', error);
    },
    clearAuth() {
      if (this.socket) {
        this.socket.close();
      }
      this.socket = null;
      this.token = '';
      this.userInfo = {};
      sessionStorage.removeItem('token');
      sessionStorage.removeItem('userInfo');
    },
    cleanUserInfo() {
      this.clearAuth();
    }
  }
})

export default pinia
