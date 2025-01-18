<template>
  <div class="modal-overlay" v-show="isVisible">
    <div class="modal-content">
      <div class="modal-header">
        <h2>聊天室</h2>
      </div>
      <div class="modal-body">
        <video autoplay playsinline class="video-player"></video>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted, reactive } from "vue";
export default {
  name: "LargeModal",
  props: {
    isVisible: false,
  },
  setup() {
    const data = reactive({
      videoPlayer: null,
    });
    onMounted(() => {
      prepareAVBase();
    });
    const prepareAVBase = async () => {
      if (!navigator.mediaDevices || !navigator.mediaDevices.enumerateDevices) {
        console.log("不支持 enumerateDevices() .");
      } else {
        navigator.mediaDevices
          .enumerateDevices()
          .then(function (devices) {
            devices.forEach(function (device) {
              console.log(
                device.kind +
                  ": " +
                  device.label +
                  " id = " +
                  device.deviceId +
                  " groupId = " +
                  device.groupId
              );
            });
          })
          .catch(function (err) {
            console.log(err.name + ": " + err.message);
          });
      }
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        console.log("getUserMedia is not supported");
      } else {
        data.videoPlayer = document.querySelector(".video-player");
        console.log(data.videoPlayer);
        var constraints = {
          video: true,
          audio: true,
        };
        try {
          var stream = await navigator.mediaDevices.getUserMedia(constraints);
          data.videoPlayer.srcObject = stream;
        } catch (error) {
          console.log(error);
        }
      }
    };
    return {
      prepareAVBase,
    };
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
  border-radius: 30px;
}

.modal-content {
  background: #fff;
  height: 500px;
  border-radius: 20px;
  width: 800px;
  box-shadow: 0 2px 15px rgb(195, 8, 8);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.video-player {
  width: 300px;
  height: 200px;
}

.modal-header {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-body {
  background-color: red;
  height: 360px;
  width: 700px;
}
</style>
