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
            <button class="icon-btn">
              <el-icon>
                <ChatRound />
              </el-icon>
            </button>
            <button class="icon-btn">
              <el-icon>
                <User />
              </el-icon>
            </button>
            <button class="icon-btn">
              <el-icon>
                <Share />
              </el-icon>
            </button>
            <button class="icon-btn">
              <el-icon>
                <Star />
              </el-icon>
            </button>
            <button class="icon-btn">
              <el-icon>
                <Search />
              </el-icon>
            </button>
          </div>
          <div class="down-bar">
            <button class="icon-btn">
              <el-icon>
                <Setting />
              </el-icon>
            </button>
            <button class="icon-btn">
              <el-icon @click="handleBackToChat"><HomeFilled /></el-icon>
            </button>
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
          <el-button class="edit-btn">编辑</el-button>
        </div>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { useRouter } from "vue-router";
export default {
  name: "OwnInfo",
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
    const handleBackToChat = () => {
      router.push("/chat");
    };
    return {
      ...toRefs(data),
      router,
      handleBackToChat,
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


</style>