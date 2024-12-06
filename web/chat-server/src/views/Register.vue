<template>
  <div class="register-wrap">
    <div class="register-window" :style="{
      boxShadow: `var(${'--el-box-shadow-dark'})`,
    }">
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
        <el-button type="primary" class="register-btn" @click="handleRegister">注册</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import axios from "axios";
export default {
  name: "register",
  setup() {
    const data = reactive({
      registerData: {
        telephone: "",
        password: "",
        nickname: "",
      },
    });
    const handleRegister = async () => {
      try {
        const response = await axios.post(
          "http://127.0.0.1:8000/register",
          data.registerData
        ); // 发送POST请求
        console.log("响应数据:", response.data); // 处理响应数据
        // 这里可以添加登录成功的逻辑，比如跳转到另一个页面或显示成功消息
      } catch (error) {
        console.error("登录失败:", error); // 处理错误
        // 这里可以添加登录失败的逻辑，比如显示错误消息
      }
    };
    return {
      ...toRefs(data),
      handleRegister,
    };
  },
};
</script>

<style>
.register-wrap {
  height: 100vh;
  background-image: url('@/assets/img/chat_server_background.jpg');
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

.register-btn, .register-btn:hover {
  background-color: rgb(229, 132, 132);
  border: none;
  color:#ffffff;
  font-weight: bold;
}
</style>