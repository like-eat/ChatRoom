<template>
  <div class="login-wrap">
    <div
      class="login-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <h2 class="login-item">登录</h2>
      <el-form
        ref="formRef"
        :model="loginData"
        label-width="70px"
        class="demo-dynamic"
      >
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
          <el-input v-model="loginData.telephone" />
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
          <el-input type="password" v-model="loginData.password" />
        </el-form-item>
      </el-form>
      <div class="login-button-container">
        <el-button type="primary" class="login-btn" @click="handleLogin"
          >登录</el-button
        >
      </div>
      <div class="go-register-button-container">
        <button class="go-register-btn" @click="handleRegister">注册</button>
      </div>
      <el-alert
        v-if="showError"
        title="Error alert"
        type="error"
        :description="errorMessage"
        show-icon
        closable
        @close="handleAlertClose"
      />
      <el-alert
        v-if="showSuccess"
        title="Success"
        type="success"
        :description="successMessage"
        show-icon
        closable
        @close="handleAlertClose"
      />
    </div>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";
export default {
  name: "Login",
  setup() {
    const data = reactive({
      loginData: {
        telephone: "",
        password: "",
      },
      showError: false,
      errorMessage: "",
      showSuccess: false,
      successMessage: "",
    });
    const router = useRouter();
    const handleLogin = async () => {
      try {
        if (data.showError || data.showSuccess) {
          console.log(
            "Alert is showing, not sending request. please close it first."
          );
          return;
        }
        if (!data.loginData.telephone || !data.loginData.password) {
          data.showError = true;
          data.errorMessage = "请填写完整登录信息。";
          return;
        }
        if (!checkTelephoneValid()) {
          data.showError = true;
          data.errorMessage = "请输入有效的手机号码。";
          return;
        }
        const response = await axios.post(
          "http://127.0.0.1:8000/login",
          data.loginData
        );
        console.log(response);
        if (response.data.code == 200) {
          data.showSuccess = true;
          data.successMessage = "登录成功！";
          console.log(response.data.message);
          sessionStorage.setItem("userInfo", response.data.data);
          router.push("/chat/sessionlist");
        } else {
          data.showError = true;
          data.errorMessage = "登录失败！请重试！";
          console.log(response.data.error);
        }
      } catch (error) {
        data.showError = true;
        data.errorMessage = "登录失败！请重试！";
      }
    };
    const handleAlertClose = () => {
      data.showError = false;
      data.errorMessage = "";
      data.showSuccess = false;
      data.successMessage = "";
    };
    const checkTelephoneValid = () => {
      const regex = /^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$/;
      return regex.test(data.loginData.telephone);
    };
    const handleRegister = () => {
      router.push("/register");
    };

    return {
      ...toRefs(data),
      router,
      handleLogin,
      handleRegister,
      handleAlertClose,
    };
  },
};
</script>

<style>
.login-wrap {
  height: 100vh;
  background-image: url("@/assets/img/chat_server_background.jpg");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.login-window {
  background-color: rgb(255, 255, 255, 0.7);
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 30px 50px;
  border-radius: 20px;
  /*opacity: 0.7;*/
}

.login-item {
  text-align: center;
  margin-bottom: 20px;
  color: #494949;
}

.login-button-container {
  display: flex;
  justify-content: center; /* 水平居中 */
  margin-top: 20px; /* 可选，根据需要调整按钮与输入框之间的间距 */
  width: 100%;
}

.login-btn,
.login-btn:hover {
  background-color: rgb(229, 132, 132);
  border: none;
  color: #ffffff;
  font-weight: bold;
}

.go-register-button-container {
  display: flex;
  flex-direction: row-reverse;
}

.go-register-btn {
  background-color: rgba(255, 255, 255, 0);
  border: none;
  cursor: pointer;
  color: #d65b54;
  font-weight: bold;
  text-decoration: underline;
  text-underline-offset: 0.2em;
}

.el-alert {
  margin-top: 20px;
}
</style>