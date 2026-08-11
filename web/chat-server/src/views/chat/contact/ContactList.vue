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
          <ContactListModal></ContactListModal>
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
import { reactive, toRefs, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAppStore } from "@/store";
import Modal from "@/components/Modal.vue";
import SmallModal from "@/components/SmallModal.vue";
import ContactListModal from "@/components/ContactListModal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
import axios from "axios";
import { ElMessage } from "element-plus";
export default {
  name: "ContactList",
  components: {
    Modal,
    SmallModal,
    ContactListModal,
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
      isCreateGroupModalVisible: false,
      isApplyContactModalVisible: false,
      isNewContactModalVisible: false,
      contactUserList: [],
      myGroupList: [],
      myJoinedGroupList: [],
      applyContactReq: {
        owner_id: "",
        contact_id: "",
        message: "",
      },
      ownListReq: {
        owner_id: "",
      },
      newContactList: [],
      applyContent: "",
    });

    const handleCreateGroup = async () => {
      try {
        data.createGroupReq.owner_id = data.userInfo.uuid;
        const response = await axios.post(
          store.backendUrl + "/group/createGroup",
          data.createGroupReq
        );
      } catch (error) {
        console.error(error);
      }
    };
    const showCreateGroupModal = () => {
      data.isCreateGroupModalVisible = true;
    };
    const quitCreateGroupModal = () => {
      data.isCreateGroupModalVisible = false;
    };
    const closeCreateGroupModal = () => {
      if (data.createGroupReq.name == "") {
        ElMessage("请输入群聊名称");
        return;
      }
      data.isCreateGroupModalVisible = false;
      handleCreateGroup();
    };
    const showApplyContactModal = () => {
      data.isApplyContactModalVisible = true;
    };
    const quitApplyContactModal = () => {
      data.isApplyContactModalVisible = false;
    };
    const closeApplyContactModal = () => {
      if (data.applyContactReq.contact_id == "") {
        ElMessage.error("请输入申请用户/群组id");
        return;
      }
      if (data.applyContactReq.contact_id[0] == 'G') {
        handleApplyGroup();
      } else {
        handleApplyContact();
      }
    };

    const showNewContactModal = () => {
      handleNewContactList();
    };

    const quitNewContactModal = () => {
      data.isNewContactModalVisible = false;
      data.newContactList = [];
    };
    const handleApplyGroup = async () => {
      try {
        data.applyContactReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.backendUrl + "/contact/applyContact",
          data.applyContactReq
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          if (rsp.data.message == "申请成功") {
            data.isApplyContactModalVisible = false;
            ElMessage.success("申请成功");
            return;
          } else {
            ElMessage.error(rsp.data.message);
          }
        } else {
          ElMessage.error("申请失败");
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleApplyContact = async () => {
      try {
        data.applyContactReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.backendUrl + "/contact/applyContact",
          data.applyContactReq
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          if (rsp.data.message == "申请成功") {
            data.isApplyContactModalVisible = false;
            ElMessage.success("申请成功");
            return;
          }
        } 
        ElMessage.error(rsp.data.message);
      } catch (error) {
        console.error(error);
      }
    };
    const handleShowUserList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const getUserListRsp = await axios.post(
          store.backendUrl + "/contact/getUserList",
          data.ownListReq
        );
        data.contactUserList = getUserListRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideUserList = () => {
      data.contactUserList = [];
    };

    const handleShowMyGroupList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const loadMyGroupRsp = await axios.post(
          store.backendUrl + "/group/loadMyGroup",
          data.ownListReq
        );
        data.myGroupList = loadMyGroupRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideMyGroupList = () => {
      data.myGroupList = [];
    };
    const handleShowMyJoinedGroupList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const loadMyJoinedGroupRsp = await axios.post(
          store.backendUrl + "/contact/loadMyJoinedGroup",
          data.ownListReq
        );
        data.myJoinedGroupList = loadMyJoinedGroupRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideMyJoinedGroupList = () => {
      data.myJoinedGroupList = [];
    };


    const handleToChatUser = async (user) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: user.user_id,
        };
        const rsp = await axios.post(
          store.backendUrl + "/session/checkOpenSessionAllowed",
          req
        );
        if (rsp.data.code == 200) {
          if (rsp.data.data == true) {
            router.push("/chat/" + user.user_id);
          } else {
            ElMessage.warning(rsp.data.message);
            console.error(rsp.data.message);
          }
        } else {
          ElMessage.error(rsp.data.message);
          console.error(rsp.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
        console.error(error);
      }
    };

    const handleToChatGroup = async (group) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: group.group_id,
        };
        const rsp = await axios.post(
          store.backendUrl + "/session/checkOpenSessionAllowed",
          req
        );
        if (rsp.data.code == 200) {
          if (rsp.data.data == true) {
            router.push("/chat/" + group.group_id);
          } else {
            ElMessage.warning(rsp.data.message);
            console.error(rsp.data.message);
          }
        } else {
          if (rsp.data.code == 400) {
            ElMessage.warning(rsp.data.message);
            console.error(rsp.data.message);
          } else {
            ElMessage.error(rsp.data.message);
            console.error(rsp.data.message);
          }
        }
      } catch (error) {
        console.error(error);
      }
    };

    const handleNewContactList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.backendUrl + "/contact/getNewContactList",
          data.ownListReq
        );
        console.log(rsp);
        data.newContactList = rsp.data.data;
        if (data.newContactList == null) {
          ElMessage.warning("没有新的好友申请");
          return;
        }
        data.isNewContactModalVisible = true;
        console.log(rsp);
      } catch (error) {
        console.error(error);
      }
    };

    const handleAgree = async (contactId) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: contactId,
        };
        const rsp = await axios.post(
          store.backendUrl + "/contact/passContactApply",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          data.newContactList = data.newContactList.filter(
            (c) => c.contact_id !== contactId
          );
        } else {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleReject = async (contactId) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: contactId,
        };
        const rsp = await axios.post(
          store.backendUrl + "/contact/refuseContactApply",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
          data.newContactList = data.newContactList.filter(
            (c) => c.contact_id !== contactId
          );
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 500) {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    }
    return {
      ...toRefs(data),
      router,
      handleCreateGroup,
      showCreateGroupModal,
      closeCreateGroupModal,
      quitCreateGroupModal,
      showApplyContactModal,
      quitApplyContactModal,
      closeApplyContactModal,
      showNewContactModal,
      quitNewContactModal,
      handleShowUserList,
      handleHideUserList,
      handleShowMyGroupList,
      handleHideMyGroupList,
      handleShowMyJoinedGroupList,
      handleHideMyJoinedGroupList,
      handleToChatUser,
      handleToChatGroup,
      handleNewContactList,
      handleAgree,
      handleReject,
    };
  },
};
</script>

<style scoped>
.contactlist-header {
  display: flex;
  flex-direction: row;
  margin-top: 10px;
  margin-bottom: 10px;
}

.contact-search-input {
  width: 185px;
  height: 30px;
  margin-left: 5px;
  margin-right: 5px;
}

.contactlist-header-right {
  width: 40px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.create-group-btn {
  background-color: rgb(252, 210.9, 210.9);
  cursor: pointer;
  border: none;
  height: 100%;
  width: 30px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 10px;
}

.create-group-icon {
  width: 15px;
  height: 15px;
}

.el-menu {
  background-color: rgb(252, 210.9, 210.9);
  width: 101%;
}

.el-menu-item {
  background-color: rgb(255, 255, 255);
  height: 45px;
}

.contactlist-user-title {
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
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.newcontact-modal-body {
  height: 70%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.newcontact-modal-footer {
  height: 10%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
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

.contactlist-avatar {
  width: 30px;
  height: 30px;
  margin-right: 20px;
}

.newcontact-list {
  width: 280px;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

.newcontact-item {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 40px;
}

.action-btn {
  background-color: rgb(252, 210.9, 210.9);
  border: none;
  cursor: pointer;
  justify-content: center;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

.contactlist-user-menu-item {
  justify-content: center;
  align-items: center;
}

.contactlist-user-item {
  width: 221px;
  height: 45px;
  display: flex;
  align-items: center;
  color: rgba(43, 42, 42, 0.893);
}

.contactlist-user-avatar {
  width: 30px;
  height: 30px;
  margin-left: 20px;
  margin-right: 20px;
}
</style>