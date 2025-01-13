<template>
  <div class="chat-wrap">
    <div
      class="chat-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <el-container class="chat-window-container">
        <el-header class="header-container">
          <svg
            t="1736782942164"
            class="header-icon"
            viewBox="0 0 1024 1024"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            p-id="2358"
            width="200"
            height="200"
          >
            <path
              d="M512 597.333333v341.333334H170.666667a341.333333 341.333333 0 0 1 341.333333-341.333334z m0-42.666666c-141.44 0-256-114.56-256-256s114.56-256 256-256 256 114.56 256 256-114.56 256-256 256z m384 170.666666h42.666667v213.333334h-341.333334v-213.333334h42.666667v-42.666666a128 128 0 1 1 256 0v42.666666z m-85.333333 0v-42.666666a42.666667 42.666667 0 1 0-85.333334 0v42.666666h85.333334z"
              p-id="2359"
              fill="#2c2c2c"
            ></path>
          </svg>
          <div class="header-title">admin管理员操作界面</div>
        </el-header>

        <el-container>
          <el-aside class="manager-aside-container">
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
          </el-aside>
          <el-main class="manager-main-container"> </el-main>
        </el-container>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted, ref } from "vue";
import { useStore } from "vuex";
import axios from "axios";
import { useRouter } from "vue-router";
import Modal from "@/components/Modal.vue";
import { checkEmailValid } from "@/assets/js/valid.js";
import { generateString } from "@/assets/js/random.js";
import SmallModal from "@/components/SmallModal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
import ContactListModal from "@/components/ContactListModal.vue";
import { ElMessage } from "element-plus";
export default {
  name: "Manager",
  components: {},
  setup() {
    const data = reactive({});
  },
};
</script>

<style scoped>
.header-container {
  height: 70px;
  border-bottom: 3px solid #ccc;
  background-color: rgb(252, 210.9, 210.9);
  border-top-right-radius: 30px;
  border-top-left-radius: 30px;
  display: flex;
  flex-direction: row;
  align-items: center;
}

.el-menu {
  width: 197px;
}

.manager-aside-container {
  width: 200px;
  border-right: 3px solid #ccc;
  display: flex;
  flex-direction: column;
  border-bottom-left-radius: 30px;
  background-color: white;
}

.header-icon {
  height: 40px;
  width: 40px;
  margin-right: 20px;
}

.header-title {
  font-family: Arial, Helvetica, sans-serif;
  font-weight: bold;
  color: rgb(37, 37, 37);
  font-size: 30px;
}

.user-manager-icon {
  height: 20px;
  width: 20px;
}

.el-menu {
}
</style>