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
    </div>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import axios from "axios";
import { useRouter } from 'vue-router';
export default {
  name: "login",
  setup() {
    const data = reactive({
      loginData: {
        telephone: "",
        password: "",
      },
    });
    const router = useRouter();
    const handleLogin = async () => {
      try {
        const response = await axios.post(
          "http://127.0.0.1:8000/login",
          data.loginData
        ); // 发送POST请求
        console.log("响应数据:", response.data); // 处理响应数据
        // 这里可以添加登录成功的逻辑，比如跳转到另一个页面或显示成功消息
      } catch (error) {
        console.error("登录失败:", error); // 处理错误
        // 这里可以添加登录失败的逻辑，比如显示错误消息
      }
    };
    const handleRegister = () => {
      router.push("/register")
    }


    return {
      ...toRefs(data),
      handleLogin,
      handleRegister,
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

.login-btn, .login-btn:hover {
  background-color: rgb(229, 132, 132);
  border: none;
  color:#ffffff;
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


</style>