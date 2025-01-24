import { createStore } from 'vuex'

export default createStore({
  state: {
    // web服务器地址
    backendUrl: 'https://192.168.2.114:8000',
    wsUrl: 'wss://192.168.2.114:8000',
    // 信令服务器地址
    // signalUrl: 'wss://127.0.0.1:8001',
    userInfo: (localStorage.getItem('userInfo') && JSON.parse(localStorage.getItem('userInfo'))) || {},
    socket: null,
  },
  getters: {
  },
  mutations: {
    setUserInfo(state, userInfo) {
      state.userInfo = userInfo;
      localStorage.setItem('userInfo', JSON.stringify(userInfo));
    },
    cleanUserInfo(state) {
      state.userInfo = {};
      localStorage.removeItem('userInfo');
    }
  },
  actions: {
  },
  modules: {
  }
})