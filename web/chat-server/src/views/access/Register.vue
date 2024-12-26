<template>
  <div class="register-wrap">
    <div
      class="register-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <h2 class="register-item">注册</h2>
      <el-form
        ref="formRef"
        :model="registerData"
        label-width="70px"
        class="demo-dynamic"
      >
        <el-form-item
          prop="nickname"
          label="昵称"
          :rules="[
            {
              required: true,
              message: '此项为必填项',
              trigger: 'blur',
            },
            {
              min: 3,
              max: 10,
              message: '昵称长度在 3 到 10 个字符',
              trigger: 'blur',
            },
          ]"
        >
          <el-input v-model="registerData.nickname" />
        </el-form-item>
        <el-form-item
          prop="telephone"
          label="账号"
          :rules="[
            {
              required: true,
              message: '此项为必填项',
              trigger: 'blur',
            },
          ]"
        >
          <el-input v-model="registerData.telephone" />
        </el-form-item>
        <el-form-item
          prop="password"
          label="密码"
          :rules="[
            {
              required: true,
              message: '此项为必填项',
              trigger: 'blur',
            },
          ]"
        >
          <el-input type="password" v-model="registerData.password" />
        </el-form-item>
      </el-form>
      <div class="register-button-container">
        <el-button type="primary" class="register-btn" @click="handleRegister"
          >注册</el-button
        >
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";
import { ElMessage } from 'element-plus';
export default {
  name: "Register",
  setup() {
    const data = reactive({
      registerData: {
        telephone: "",
        password: "",
        nickname: "",
      },
    });
    const router = useRouter();
    const handleRegister = async () => {
      try {
        if (
          !data.registerData.nickname ||
          !data.registerData.telephone ||
          !data.registerData.password
        ) {
          ElMessage.error("请填写完整注册信息。");
          return;
        }
        if (
          data.registerData.nickname.length < 3 ||
          data.registerData.nickname.length > 10
        ) {
          ElMessage.error("昵称长度在 3 到 10 个字符。");
          return;
        }
        if (!checkTelephoneValid()) {
          ElMessage.error("请输入有效的手机号码。");
          return;
        }
        const response = await axios.post(
          "http://127.0.0.1:8000/register",
          data.registerData
        ); // 发送POST请求
        if (response.data.code == 200) {
          ElMessage.success(response.data.message);
          console.log(response.data.message);
          sessionStorage.setItem("userInfo", response.data.data);
          router.push("/chat/sessionlist");
        } else {
          ElMessage.error(response.data.message);
          console.log(response.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
        console.log(error);
      }
    };
    const checkTelephoneValid = () => {
      const regex = /^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$/;
      return regex.test(data.registerData.telephone);
    };
    return {
      ...toRefs(data),
      router,
      handleRegister,
    };
  },
};
</script>

<style>
.register-wrap {
  height: 100vh;
  background-image: url("@/assets/img/chat_server_background.jpg");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.register-window {
  background-color: rgb(255, 255, 255, 0.7);
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 30px 50px;
  border-radius: 20px;
}

.register-item {
  text-align: center;
  margin-bottom: 20px;
  color: #494949;
}

.register-button-container {
  display: flex;
  justify-content: center; /* 水平居中 */
  margin-top: 20px; /* 可选，根据需要调整按钮与输入框之间的间距 */
  width: 100%;
}

.register-btn,
.register-btn:hover {
  background-color: rgb(229, 132, 132);
  border: none;
  color: #ffffff;
  font-weight: bold;
}

.el-alert {
  margin-top: 20px;
}
</style>