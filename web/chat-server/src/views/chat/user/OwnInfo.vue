<template>
  <div class="chat-wrap">
    <div
      class="chat-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <el-container class="chat-window-container">
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
                <button class="icon-btn">
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
                <button class="icon-btn" @click="handleToCantactList">
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
        <div class="owner-info-window">
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
          <el-button class="edit-btn" @click="showModal">编辑</el-button>
        </div>
        <Modal :isVisible="isModalVisible">
          <template v-slot:header>
            <div class="modal-header">
              <div class="modal-quit-btn-container">
                <button class="modal-quit-btn" @click="quitModal">
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
              <el-button
                class="modal-close-btn"
                @click="closeModal"
              >
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
import { useRouter } from "vue-router";
import Modal from "@/components/Modal.vue";
import { checkEmailValid } from "@/assets/js/valid.js";
export default {
  name: "OwnInfo",
  components: {
    Modal,
  },
  setup() {
    const data = reactive({
      userInfo: {
        uuid: "",
        nickname: "",
        telephone: "",
        avatar: "",
        email: "",
        gender: "",
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
      isModalVisible: false,
    });
    const router = useRouter();
    onMounted(() => {
      const userInfoStr = sessionStorage.getItem("userInfo");
      if (userInfoStr) {
        try {
          data.userInfo = JSON.parse(userInfoStr);
          if (data.userInfo.gender == false) {
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
    const showModal = () => {
      data.isModalVisible = true;
    };
    const closeModal = () => {
      if (
        data.updateInfo.nickname == "" &&
        data.updateInfo.avatar == "" &&
        data.updateInfo.email == "" &&
        data.updateInfo.birthday == "" &&
        data.updateInfo.signature == ""
      ) {
        alert("请至少修改一项");
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
          alert("请输入有效的邮箱。");
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
      data.isModalVisible = false;
    };
    const quitModal = () => {
      data.isModalVisible = false;
    };

    const handleToCantactList = () => {
      router.push("/chat/contactlist");
    }
    return {
      ...toRefs(data),
      router,
      showModal,
      closeModal,
      quitModal,
      handleToCantactList,
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
  width: 400px;
}

.modal-footer {
  height: 25%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

</style>