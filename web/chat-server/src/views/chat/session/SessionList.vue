<template>
  <div class="chat-wrap">
    <div
      class="chat-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <el-container class="chat-window-container">
        <el-aside class="aside-container">
          <NavigationModal></NavigationModal>
          <div class="sessionlist-container">
            <div class="sessionlist-header">
              <el-input
                v-model="contactSearch"
                class="contact-search-input"
                placeholder="搜索会话"
                size="small"
                suffix-icon="Search"
              />
            </div>
            <div class="contactlist-body">
              <div class="contactlist-user">
                <el-menu
                  router
                  unique-opened
                  @open="handleShowUserSessionList"
                  @close="handleHideUserSessionList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="sessionlist-title">用户</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="user in userSessionList"
                    :key="user.user_id"
                    @click="handleToChatUser(user)"
                  >
                    <img :src="user.avatar" class="sessionlist-avatar" />
                    {{ user.user_name }}
                  </el-menu-item>
                </el-menu>
                <el-menu
                  router
                  unique-opened
                  @open="handleShowGroupSessionList"
                  @close="handleHideGroupSessionList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="sessionlist-title">群聊</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="group in groupSessionList"
                    :key="group.group_id"
                    @click="handleToChatGroup(group)"
                  >
                    <img :src="group.avatar" class="sessionlist-avatar" />
                    {{ group.group_name }}
                  </el-menu-item>
                </el-menu>
              </div>
            </div>
          </div>
        </el-aside>
        <el-container class="chat-container">
          <el-header>
            <h2 class="chat-name"></h2>
          </el-header>
          <el-main class="main-container">
            <div class="chat-screen"></div>
            <div class="tool-bar">
              <div class="tool-bar-left">
                <el-tooltip
                  effect="customized"
                  content="聊天记录"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733504061769"
                      class="record-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="5492"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M695.04 194.32H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h596.96c18.32 0 33.16 14.85 33.16 33.16 0 18.31-14.84 33.16-33.16 33.16zM298.97 393.3H96.19c-17.27 0-31.27-14-31.27-31.27v-3.79c0-17.27 14-31.27 31.27-31.27h202.78c17.27 0 31.27 14 31.27 31.27v3.79c-0.01 17.28-14.01 31.27-31.27 31.27zM230.74 592.29H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h132.66c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16zM230.74 791.28H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h132.66c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16zM728.2 691.78H595.55c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16H728.2c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16z"
                        p-id="5493"
                      ></path>
                      <path
                        d="M562.38 658.62V525.96c0-18.32 14.85-33.16 33.16-33.16 18.32 0 33.16 14.85 33.16 33.16v132.66c0 18.32-14.85 33.16-33.16 33.16-18.31 0-33.16-14.85-33.16-33.16z"
                        p-id="5494"
                      ></path>
                      <path
                        d="M960.35 625.45c0 183.16-148.48 331.64-331.64 331.64S297.07 808.62 297.07 625.45s148.48-331.64 331.64-331.64 331.64 148.48 331.64 331.64zM628.71 360.14c-146.53 0-265.31 118.79-265.31 265.31s118.79 265.31 265.31 265.31 265.31-118.79 265.31-265.31-118.78-265.31-265.31-265.31z"
                        p-id="5495"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="全文复制"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503137487"
                      class="copy-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="3442"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M394.666667 106.666667h448a74.666667 74.666667 0 0 1 74.666666 74.666666v448a74.666667 74.666667 0 0 1-74.666666 74.666667H394.666667a74.666667 74.666667 0 0 1-74.666667-74.666667V181.333333a74.666667 74.666667 0 0 1 74.666667-74.666666z m0 64a10.666667 10.666667 0 0 0-10.666667 10.666666v448a10.666667 10.666667 0 0 0 10.666667 10.666667h448a10.666667 10.666667 0 0 0 10.666666-10.666667V181.333333a10.666667 10.666667 0 0 0-10.666666-10.666666H394.666667z m245.333333 597.333333a32 32 0 0 1 64 0v74.666667a74.666667 74.666667 0 0 1-74.666667 74.666666H181.333333a74.666667 74.666667 0 0 1-74.666666-74.666666V394.666667a74.666667 74.666667 0 0 1 74.666666-74.666667h74.666667a32 32 0 0 1 0 64h-74.666667a10.666667 10.666667 0 0 0-10.666666 10.666667v448a10.666667 10.666667 0 0 0 10.666666 10.666666h448a10.666667 10.666667 0 0 0 10.666667-10.666666v-74.666667z"
                        fill="#000000"
                        p-id="3443"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
              </div>
            </div>
          </el-main>
          <el-footer>
            <div class="chat-input">
              <el-input
                v-model="chatMessage"
                type="textarea"
                show-word-limit
                maxlength="500"
                :autosize="{ minRows: 7.9, maxRows: 7 }"
                placeholder="请输入内容"
              />
            </div>
            <div class="chat-send">
              <el-button class="send-btn">发送</el-button>
            </div>
          </el-footer>
        </el-container>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted, ref } from "vue";
import { onBeforeRouteUpdate, useRouter } from "vue-router";
import { ElMessageBox, ElMessage } from "element-plus";
import { useAppStore } from "@/store";
import axios from "axios";
import Modal from "@/components/Modal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
export default {
  name: "ContactList",
  components: {
    Modal,
    NavigationModal,
  },

  setup() {
    const router = useRouter();
    const store = useAppStore();
    const data = reactive({
      chatMessage: "",
      chatName: "",
      userInfo: store.userInfo,
      contactSearch: "",
      createGroupReq: {
        owner_id: "",
        name: "",
        notice: "",
        avatar: "",
      },
      ownListReq: {
        owner_id: "",
      },
      userSessionList: [],
      groupSessionList: [],
    });
    onMounted(() => {
      // if (data.userInfo.avatar.startswith("/static")) {
      //   data.userInfo.avatar = data.
      // }
    });
    const handleToChatUser = (user) => {
      router.push("/chat/" + user.user_id);
    };
    const handleToChatGroup = (group) => {
      router.push("/chat/" + group.group_id);
    };

    const handleShowUserSessionList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const userSessionListRsp = await axios.post(
          store.backendUrl + "/session/getUserSessionList",
          data.ownListReq
        );
        if (userSessionListRsp.data.data) {
          for (let i = 0; i < userSessionListRsp.data.data.length; i++) {
            if (!userSessionListRsp.data.data[i].avatar.startsWith("http")) {
              userSessionListRsp.data.data[i].avatar =
                store.backendUrl + userSessionListRsp.data.data[i].avatar;
            }
          }
        }
        data.userSessionList = userSessionListRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideUserSessionList = () => {
      data.userSessionList = [];
    };
    const handleShowGroupSessionList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const groupSessionListRsp = await axios.post(
          store.backendUrl + "/session/getGroupSessionList",
          data.ownListReq
        );
        if (groupSessionListRsp.data.data) {
          for (let i = 0; i < groupSessionListRsp.data.data.length; i++) {
            if (!groupSessionListRsp.data.data[i].avatar.startsWith("http")) {
              groupSessionListRsp.data.data[i].avatar =
                store.backendUrl + groupSessionListRsp.data.data[i].avatar;
            }
          }
        }
        data.groupSessionList = groupSessionListRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideGroupSessionList = () => {
      data.groupSessionList = [];
    };
    const handleContextMenu = (event, group) => {
      event.preventDefault(); // 阻止默认的右键菜单
      // 显示自定义的删除选项
      ElMessageBox.confirm("确定要删除该会话组吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          // 执行删除操作
          this.deleteGroup(group);
        })
        .catch(() => {
          // 取消删除操作
        });
    };
    return {
      ...toRefs(data),
      router,
      handleToChatUser,
      handleToChatGroup,
      handleShowUserSessionList,
      handleHideUserSessionList,
      handleShowGroupSessionList,
      handleHideGroupSessionList,
      handleContextMenu,
    };
  },
};
</script>

<style scoped>
.sessionlist-header {
  display: flex;
  flex-direction: row;
  width: 100%;
  margin-top: 10px;
  margin-bottom: 10px;
}

.contact-search-input {
  width: 215px;
  height: 30px;
  margin-left: 5px;
  margin-right: 2px;
}

.el-menu {
  background-color: rgb(252, 210.9, 210.9);
  width: 100%;
}

.el-menu-item {
  background-color: rgb(255, 255, 255);
  height: 45px;
}

.sessionlist-title {
  font-family: Arial, Helvetica, sans-serif;
}

h3 {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(69, 69, 68);
}

.modal-quit-btn-container {
  height: 30%;
  width: 100%;
  display: flex;
  flex-direction: row-reverse;
}

.modal-quit-btn {
  background-color: rgba(255, 255, 255, 0);
  color: rgb(229, 25, 25);
  padding: 15px;
  border: none;
  cursor: pointer;
  position: fixed;
  justify-content: center;
  align-items: center;
}

.modal-header {
  height: 20%;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  /*background-color:aqua;*/
}

.modal-body {
  height: 55%;
  width: 400px;
}

.modal-footer {
  height: 25%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-header-title {
  height: 70%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.sessionlist-avatar {
  width: 30px;
  height: 30px;
  margin-right: 20px;
}
</style>