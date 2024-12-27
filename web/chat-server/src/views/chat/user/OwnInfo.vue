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
          <div class="navigation-bar">
            <div class="up-bar">
              <button class="avatar-btn">
                <el-avatar
                  src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
                />
              </button>
            </div>
            <div class="middle-bar">
              <el-tooltip
                effect="customized"
                content="会话聊天"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn" @click="handleToSessionList">
                  <el-icon>
                    <ChatRound />
                  </el-icon>
                </button>
              </el-tooltip>
              <el-tooltip
                effect="customized"
                content="通讯录管理"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn" @click="handleToContactList">
                  <el-icon>
                    <User />
                  </el-icon>
                </button>
              </el-tooltip>
              <el-tooltip
                effect="customized"
                content="朋友圈"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn">
                  <el-icon>
                    <Share />
                  </el-icon>
                </button>
              </el-tooltip>
              <el-tooltip
                effect="customized"
                content="我的收藏"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn">
                  <el-icon>
                    <Star />
                  </el-icon>
                </button>
              </el-tooltip>
              <el-tooltip
                effect="customized"
                content="搜索"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn">
                  <el-icon>
                    <Search />
                  </el-icon>
                </button>
              </el-tooltip>
            </div>
            <div class="down-bar">
              <el-tooltip
                effect="customized"
                content="设置"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn">
                  <el-icon>
                    <Setting />
                  </el-icon>
                </button>
              </el-tooltip>
              <el-tooltip
                effect="customized"
                content="我的主页"
                placement="left"
                hide-after="0"
                enterable="false"
              >
                <button class="icon-btn">
                  <el-icon><HomeFilled /></el-icon>
                </button>
              </el-tooltip>
            </div>
          </div>
          <div class="contactlist-container">
            <div class="contactlist-header">
              <el-input
                v-model="contactSearch"
                class="contact-search-input"
                placeholder="搜索联系人/群聊"
                size="small"
                suffix-icon="Search"
              />
              <div class="contactlist-header-right">
                <el-dropdown placement="bottom" trigger="click">
                  <button class="create-group-btn">
                    <svg
                      t="1733664667695"
                      class="create-group-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="2875"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M488.021333 96a248.021333 248.021333 0 1 1-17.92 495.36l-1.749333 0.341333-4.352 0.298667A304 304 0 0 0 160 896a32 32 0 1 1-64 0 368.170667 368.170667 0 0 1 250.026667-348.672A247.978667 247.978667 0 0 1 488.021333 96z m288 528a32 32 0 0 1 32 32l-0.042666 87.978667H896a32 32 0 0 1 31.701333 27.690666l0.298667 4.352a32 32 0 0 1-32 32l-88.021333-0.042666V896a32 32 0 0 1-27.648 31.701333l-4.352 0.298667a32 32 0 0 1-32-32v-88.021333h-87.978667a32 32 0 0 1-31.701333-27.648l-0.298667-4.352a32 32 0 0 1 32-32h87.978667v-87.978667a32 32 0 0 1 27.690666-31.701333zM488.021333 160a184.021333 184.021333 0 1 0 0 368 184.021333 184.021333 0 0 0 0-368z"
                        fill="#2c2c2c"
                        p-id="2876"
                      ></path>
                    </svg>
                  </button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="showCreateGroupModal">
                        创建群聊
                      </el-dropdown-item>
                      <el-dropdown-item @click="showApplyContactModal">
                        添加用户/群聊
                      </el-dropdown-item>
                      <el-dropdown-item @click="showNewContactModal">
                        新的朋友
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                <SmallModal :isVisible="isNewContactModalVisible">
                  <template v-slot:header>
                    <div class="modal-header">
                      <div class="modal-quit-btn-container">
                        <button
                          class="modal-quit-btn"
                          @click="quitNewContactModal"
                        >
                          <el-icon><Close /></el-icon>
                        </button>
                      </div>
                      <div class="modal-header-title">
                        <h3>新的朋友</h3>
                      </div>
                    </div>
                  </template>
                  <template v-slot:body>
                    <div class="newcontact-modal-body">
                      <el-scrollbar max-height="400px">
                        <ul
                          class="newcontact-list"
                          style="list-style-type: none"
                        >
                          <li
                            v-for="newContact in newContactList"
                            :key="newContact.contact_id"
                            class="newcontact-item"
                          >
                            <div
                              style="
                                display: flex;
                                align-items: center;
                                justify-content: center;
                              "
                            >
                              <img
                                :src="newContact.contact_avatar"
                                style="
                                  width: 30px;
                                  height: 30px;
                                  margin-right: 10px;
                                "
                              />

                              <el-tooltip
                                effect="customized"
                                :content="newContact.message"
                                placement="top"
                                hide-after="0"
                                enterable="false"
                              >
                                <div>
                                  {{ newContact.contact_name }}
                                </div>
                              </el-tooltip>
                            </div>
                            <el-dropdown placement="right" trigger="click">
                              <el-button class="action-btn"> 去处理 </el-button>
                              <template #dropdown>
                                <el-dropdown-menu>
                                  <el-dropdown-item
                                    @click="handleAgree(newContact.contact_id)"
                                    >同意</el-dropdown-item
                                  >
                                  <el-dropdown-item
                                    @click="handleReject(newContact.contact_id)"
                                  >
                                    拒绝
                                  </el-dropdown-item>
                                  <el-dropdown-item
                                    @click="handleBlack(newContact.contact_id)"
                                  >
                                    拉黑
                                  </el-dropdown-item>
                                </el-dropdown-menu>
                              </template>
                            </el-dropdown>
                          </li>
                        </ul>
                      </el-scrollbar>
                    </div>
                  </template>
                  <template v-slot:footer>
                    <div class="newcontact-modal-footer"></div>
                  </template>
                </SmallModal>
                <SmallModal :isVisible="isApplyContactModalVisible">
                  <template v-slot:header>
                    <div class="modal-header">
                      <div class="modal-quit-btn-container">
                        <button
                          class="modal-quit-btn"
                          @click="quitApplyContactModal"
                        >
                          <el-icon><Close /></el-icon>
                        </button>
                      </div>
                      <div class="modal-header-title">
                        <h3>添加用户/群聊</h3>
                      </div>
                    </div>
                  </template>
                  <template v-slot:body>
                    <div class="modal-body">
                      <el-form
                        ref="formRef"
                        :model="applyContactReq"
                        label-width="100px"
                        class="apply-contact-form"
                      >
                        <el-form-item
                          prop="name"
                          label="用户/群聊id"
                          :rules="[
                            {
                              required: true,
                              message: '此项为必填项',
                              trigger: 'blur',
                            },
                          ]"
                        >
                          <el-input
                            v-model="applyContactReq.contact_id"
                            placeholder="请填写申请的用户/群聊id"
                          />
                        </el-form-item>
                        <el-form-item prop="message" label="申请消息">
                          <el-input
                            v-model="applyContactReq.message"
                            placeholder="选填，填写更容易通过"
                            type="textarea"
                            show-word-limit
                            maxlength="100"
                            :autosize="{ minRows: 3, maxRows: 3 }"
                          />
                        </el-form-item>
                      </el-form>
                    </div>
                  </template>
                  <template v-slot:footer>
                    <div class="modal-footer">
                      <el-button
                        class="modal-close-btn"
                        @click="closeApplyContactModal"
                      >
                        完成
                      </el-button>
                    </div>
                  </template>
                </SmallModal>
                <Modal :isVisible="isCreateGroupModalVisible">
                  <template v-slot:header>
                    <div class="modal-header">
                      <div class="modal-quit-btn-container">
                        <button
                          class="modal-quit-btn"
                          @click="quitCreateGroupModal"
                        >
                          <el-icon><Close /></el-icon>
                        </button>
                      </div>
                      <div class="modal-header-title">
                        <h3>创建群聊</h3>
                      </div>
                    </div>
                  </template>
                  <template v-slot:body>
                    <div class="modal-body">
                      <el-form
                        ref="formRef"
                        :model="createGroupReq"
                        label-width="80px"
                        class="demo-dynamic"
                      >
                        <el-form-item
                          prop="name"
                          label="群名称"
                          :rules="[
                            {
                              required: true,
                              message: '此项为必填项',
                              trigger: 'blur',
                            },
                          ]"
                        >
                          <el-input
                            v-model="createGroupReq.name"
                            placeholder="必填"
                          />
                        </el-form-item>
                        <el-form-item prop="notice" label="群公告">
                          <el-input
                            v-model="createGroupReq.notice"
                            type="textarea"
                            show-word-limit
                            maxlength="500"
                            :autosize="{ minRows: 3, maxRows: 3 }"
                            placeholder="选填"
                          />
                        </el-form-item>
                        <el-form-item
                          prop="add_mode"
                          label="加群方式"
                          :rules="[
                            {
                              required: true,
                              message: 'Please select activity resource',
                              trigger: 'change',
                            },
                          ]"
                        >
                          <el-radio-group v-model="createGroupReq.add_mode">
                            <el-radio :value="0">直接加入</el-radio>
                            <el-radio :value="1">群主审核</el-radio>
                          </el-radio-group>
                        </el-form-item>
                        <el-form-item prop="avatar" label="群头像">
                          <el-input
                            v-model="createGroupReq.avatar"
                            placeholder="选填"
                          />
                        </el-form-item>
                      </el-form>
                    </div>
                  </template>
                  <template v-slot:footer>
                    <div class="modal-footer">
                      <el-button
                        class="modal-close-btn"
                        @click="closeCreateGroupModal"
                      >
                        完成
                      </el-button>
                    </div>
                  </template>
                </Modal>
              </div>
            </div>
            <div class="contactlist-body">
              <div class="contactlist-user">
                <el-menu
                  router
                  unique-opened
                  @open="handleShowUserList"
                  @close="handleHideUserList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="contactlist-user-title">联系人</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="user in contactUserList"
                    :key="user.user_id"
                    @click="handleToChatUser(user)"
                  >
                    <img :src="user.avatar" class="contactlist-avatar" />
                    {{ user.user_name }}
                  </el-menu-item>
                </el-menu>
                <el-menu
                  router
                  unique-opened
                  @open="handleShowMyGroupList"
                  @close="handleHideMyGroupList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="contactlist-user-title">我创建的群聊</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="group in myGroupList"
                    :key="group.group_id"
                    @click="handleToChatGroup(group)"
                  >
                    <img :src="group.avatar" class="contactlist-avatar" />
                    {{ group.group_name }}
                  </el-menu-item>
                </el-menu>
                <el-menu
                  router
                  unique-opened
                  @open="handleShowMyJoinedGroupList"
                  @close="handleHideMyJoinedGroupList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="contactlist-user-title">我加入的群聊</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="group in myJoinedGroupList"
                    :key="group.group_id"
                    @click="handleToChatGroup(group)"
                  >
                    <img :src="group.avatar" class="contactlist-avatar" />
                    {{ group.group_name }}
                  </el-menu-item>
                </el-menu>
              </div>
            </div>
          </div>
        </el-aside>
        <div class="owner-info-window">
          <div class="my-homepage-title"><h2>我的主页</h2></div>

          <p class="owner-prefix">用户id：{{ userInfo.uuid }}</p>
          <p class="owner-prefix">昵称：{{ userInfo.nickname }}</p>
          <p class="owner-prefix">电话：{{ userInfo.telephone }}</p>
          <p class="owner-prefix">邮箱：{{ userInfo.email }}</p>
          <p class="owner-prefix">性别：{{ userInfo.gender }}</p>
          <p class="owner-prefix">生日：{{ userInfo.birthday }}</p>
          <p class="owner-prefix">个性签名：{{ userInfo.signature }}</p>
          <p class="owner-prefix">
            加入kama chat server的时间：{{ userInfo.created_at }}
          </p>
          <div class="owner-opt">
            <p class="owner-prefix">头像：</p>
            <img style="width: 40px; height: 40px" :src="userInfo.avatar" />
          </div>
        </div>
        <div class="edit-window">
          <el-button class="edit-btn" @click="showMyInfoModal">编辑</el-button>
        </div>
        <Modal :isVisible="isMyInfoModalVisible">
          <template v-slot:header>
            <div class="modal-header">
              <div class="modal-quit-btn-container">
                <button class="modal-quit-btn" @click="quitMyInfoModal">
                  <el-icon><Close /></el-icon>
                </button>
              </div>
              <div class="modal-header-title">
                <h3>修改主页</h3>
              </div>
            </div>
          </template>
          <template v-slot:body>
            <div class="modal-body">
              <el-form
                ref="formRef"
                :model="updateInfo"
                label-width="70px"
                class="demo-dynamic"
              >
                <el-form-item
                  prop="nickname"
                  label="昵称"
                  :rules="[
                    {
                      min: 3,
                      max: 10,
                      message: '昵称长度在 3 到 10 个字符',
                      trigger: 'blur',
                    },
                  ]"
                >
                  <el-input v-model="updateInfo.nickname" placeholder="选填" />
                </el-form-item>
                <el-form-item prop="email" label="邮箱">
                  <el-input v-model="updateInfo.email" placeholder="选填" />
                </el-form-item>
                <el-form-item prop="birthday" label="生日">
                  <el-input
                    v-model="updateInfo.birthday"
                    placeholder="选填，格式为2024.1.1"
                  />
                </el-form-item>
                <el-form-item prop="signature" label="个性签名">
                  <el-input v-model="updateInfo.signature" placeholder="选填" />
                </el-form-item>
                <el-form-item prop="avatar" label="头像">
                  <el-input v-model="updateInfo.avatar" placeholder="选填" />
                </el-form-item>
              </el-form>
            </div>
          </template>
          <template v-slot:footer>
            <div class="modal-footer">
              <el-button class="modal-close-btn" @click="closeMyInfoModal">
                完成
              </el-button>
            </div>
          </template>
        </Modal>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { useStore } from "vuex";
import axios from "axios";
import { useRouter } from "vue-router";
import Modal from "@/components/Modal.vue";
import { checkEmailValid } from "@/assets/js/valid.js";
import SmallModal from "@/components/SmallModal.vue";
import { ElMessage } from "element-plus";
export default {
  name: "OwnInfo",
  components: {
    Modal,
    SmallModal,
  },
  setup() {
    const data = reactive({
      userInfo: {
        uuid: "",
        nickname: "",
        telephone: "",
        avatar: "",
        email: "",
        gender: null,
        birthday: "",
        signature: "",
        created_at: "",
      },
      updateInfo: {
        nickname: "",
        avatar: "",
        email: "",
        birthday: "",
        signature: "",
      },
      isMyInfoModalVisible: false,
      isCreateGroupModalVisible: false,
      isNewContactModalVisible: false,
      isApplyJoinGroupModalVisible: false,
      contactSearch: "",
      createGroupReq: {
        owner_id: "",
        name: "",
        notice: "",
        add_mode: null,
        avatar: "",
      },
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
    const router = useRouter();
    const store = useStore();
    onMounted(() => {
      const userInfoStr = sessionStorage.getItem("userInfo");
      if (userInfoStr) {
        try {
          data.userInfo = JSON.parse(userInfoStr);
          if (data.userInfo.gender == 0) {
            data.userInfo.gender = "男";
          } else {
            data.userInfo.gender = "女";
          }
          console.log("反序列化用户信息成功:", data.userInfo);
        } catch (error) {
          console.error("反序列化用户信息时出错:", error);
          data.userInfo = {};
        }
      } else {
        data.userInfo = {};
      }
    });
    const showMyInfoModal = () => {
      data.isMyInfoModalVisible = true;
    };
    const closeMyInfoModal = () => {
      if (
        data.updateInfo.nickname == "" &&
        data.updateInfo.avatar == "" &&
        data.updateInfo.email == "" &&
        data.updateInfo.birthday == "" &&
        data.updateInfo.signature == ""
      ) {
        ElMessage("请至少修改一项");
        return;
      }
      if (data.updateInfo.nickname != "") {
        if (
          data.updateInfo.nickname.length < 3 ||
          data.updateInfo.nickname.length > 10
        ) {
          return;
        }
      }
      if (data.updateInfo.email != "") {
        if (!checkEmailValid(data.updateInfo.email)) {
          ElMessage("请输入有效的邮箱。");
          return;
        }
      }
      if (data.updateInfo.nickname != "") {
        data.userInfo.nickname = data.updateInfo.nickname;
      }
      if (data.updateInfo.email != "") {
        data.userInfo.email = data.updateInfo.email;
      }
      if (data.updateInfo.avatar != "") {
        data.userInfo.avatar = data.updateInfo.avatar;
      }

      if (data.updateInfo.birthday != "") {
        data.userInfo.birthday = data.updateInfo.birthday;
      }
      if (data.updateInfo.signature != "") {
        data.userInfo.signature = data.updateInfo.signature;
      }
      const userInfoStr = JSON.stringify(data.userInfo);
      localStorage.setItem("userInfo", userInfoStr);
      data.isMyInfoModalVisible = false;
    };
    const quitMyInfoModal = () => {
      data.isMyInfoModalVisible = false;
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
      if (data.createGroupReq.add_mode == null) {
        ElMessage("请选择加群方式");
        return;
      }
      data.isCreateGroupModalVisible = false;
      handleCreateGroup();
    };
    const handleToContactList = () => {
      router.push("/chat/contactlist");
    };

    const handleToSessionList = () => {
      router.push("/chat/sessionlist");
    };

    const handleCreateGroup = async () => {
      try {
        data.createGroupReq.owner_id = data.userInfo.uuid;
        const response = await axios.post(
          store.state.backendUrl + "/group/createGroup",
          data.createGroupReq
        );
      } catch (error) {
        console.error(error);
      }
    };

    const handleShowUserList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const getUserListRsp = await axios.post(
          store.state.backendUrl + "/contact/getUserList",
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
          store.state.backendUrl + "/group/loadMyGroup",
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
          store.state.backendUrl + "/contact/loadMyJoinedGroup",
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
      handleApplyContact();
    };

    const showNewContactModal = () => {
      handleNewContactList();
    };

    const quitNewContactModal = () => {
      data.isNewContactModalVisible = false;
      data.newContactList = [];
    };

    const handleApplyContact = async () => {
      try {
        data.applyContactReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/applyContact",
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

    const handleToChatUser = async (user) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: user.user_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/session/checkOpenSessionAllowed",
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
      router.push("/chat/" + group.group_id);
    };

    const handleNewContactList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/getNewContactList",
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
          store.state.backendUrl + "/contact/passContactApply",
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

    const handleCancelBlack = async (user) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: user.user_id,
        };
        const rsp = await axios.post(store.state.backendUrl + "/contact/cancelBlackContact", req);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 500) {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
        console.error(error);
      }
    };
    return {
      ...toRefs(data),
      router,
      showMyInfoModal,
      closeMyInfoModal,
      quitMyInfoModal,
      showCreateGroupModal,
      closeCreateGroupModal,
      quitCreateGroupModal,
      handleToContactList,
      handleToSessionList,
      handleCreateGroup,
      handleShowUserList,
      handleHideUserList,
      handleShowMyGroupList,
      handleHideMyGroupList,
      handleShowMyJoinedGroupList,
      handleHideMyJoinedGroupList,
      showApplyContactModal,
      quitApplyContactModal,
      closeApplyContactModal,
      showNewContactModal,
      quitNewContactModal,
      handleToChatUser,
      handleToChatGroup,
      handleNewContactList,
      handleAgree,
      handleCancelBlack,
    };
  },
};
</script>

<style scoped>
.owner-info-window {
  width: 84%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.owner-prefix {
  font-family: Arial, Helvetica, sans-serif;
  margin: 6px;
}

.owner-opt {
  margin: 6px;
  display: flex;
  flex-direction: row;
}

.edit-window {
  width: 10%;
  display: flex;
  flex-direction: column-reverse;
}

.modal-header-title {
  height: 70%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
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

h2 {
  margin-bottom: 20px;
  font-family: Arial, Helvetica, sans-serif;
}

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
</style>