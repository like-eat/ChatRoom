<template>
  <div class="modal-overlay" v-show="isVisible">
    <div class="modal-content">
      <div class="modal-header">
        <h2>聊天室</h2>
      </div>
      <div class="modal-body">
        <video autoplay playsinline class="video-player"></video>
      </div>
      <div class="modal-footer">
      <el-button class="modal-footer-btn" @click="emitSession">发起通话</el-button>
      <el-button class="modal-footer-btn">挂断通话</el-button>
      <el-button class="modal-footer-btn">退出聊天室</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted, reactive } from "vue";
import { useStore } from "vuex";
import { ElMessage } from 'element-plus';
export default {
  name: "LargeModal",
  props: {
    isVisible: false,
    ownerId: null,
  },
  setup(props) {
    const store = useStore();
    const data = reactive({
      videoPlayer: null,
      rtcPeerConn: null,
      ICE_CFG: {},
      contactId: null,
      wsConn: null,
      localStream: null,
      remoteStream: null,
      remoteVideo: document.querySelector("video.remote-video"),
      localVideo: document.querySelector("video.local-video"),
      curContactList: [],
    });
    onMounted(() => { 
      // prepareAVBase();
    });
    const emitSession = () => {
      login();
    }
    // const prepareAVBase = async () => {
    //   if (!navigator.mediaDevices || !navigator.mediaDevices.enumerateDevices) {
    //     console.log("不支持 enumerateDevices() .");
    //   } else {
    //     navigator.mediaDevices
    //       .enumerateDevices()
    //       .then(function (devices) {
    //         devices.forEach(function (device) {
    //           console.log(
    //             device.kind +
    //               ": " +
    //               device.label +
    //               " id = " +
    //               device.deviceId +
    //               " groupId = " +
    //               device.groupId
    //           );
    //         });
    //       })
    //       .catch(function (err) {
    //         console.log(err.name + ": " + err.message);
    //       });
    //   }
    //   if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
    //     console.log("getUserMedia is not supported");
    //   } else {
    //     data.videoPlayer = document.querySelector(".video-player");
    //     console.log(data.videoPlayer);
    //     var constraints = {
    //       video: true,
    //       audio: true,
    //     };
    //     try {
    //       var stream = await navigator.mediaDevices.getUserMedia(constraints);
    //       data.videoPlayer.srcObject = stream;
    //     } catch (error) {
    //       console.log(error);
    //     }
    //   }
    // };
    const createRtcPeerConnection = () => {
      if (data.rtcPeerConn) {
        console.log("peer connection has already been created.");
        return;
      }
      rtcPeerConn = new RTCPeerConnection(data.ICE_CFG);
      rtcPeerConn.onicecandidate = (event) => {
        if (event.candidate) {
          var proxyCandidateMessage = {
            messageId: "PROXY",
            type: "candidate",
            fromId: props.ownerId,
            toId: data.contactId,
            messageData: {
              candidate: event.candidate,
            },
          };
          data.wsConn.send(JSON.stringify(proxyCandidateMessage));
        }
      };
      rtcPeerConn.oniceconnectionstatechange = (event) => {
        console.log(
          "oniceconnectionstatechange",
          rtcPeerConn.iceConnectionState
        );
      };
      // 对端传来媒体轨道
      rtcPeerConn.ontrack = (event) => {
        if (remoteStream === null) {
          data.remoteStream = new MediaStream();
          data.remoteVideo.srcObject = data.remoteStream;
          data.remoteVideo.style.display = "inline-block";
        }
        data.remoteStream.addTrack(event.track);
      };
    };

    const closeRtcPeerConnection = () => {
      if (data.rtcPeerConn) {
        data.rtcPeerConn.close();
        data.rtcPeerConn = null;
      }
    };

    const getLocalMediaStream = () => {
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        console.log("getUserMedia is not supported!");
        return null;
      } else if (data.localStream) {
        console.log("localStream already exist.");
        return data.localStream;
      } else {
        var constraints = {
          video: true,
          audio: true,
        };
        return navigator.mediaDevices.getUserMedia(constraints);
      }
    };

    const closeLocalMediaStream = () => {
      if (data.localStream != null) {
        data.localStream.getTracks().forEach((track) => {
          track.stop();
        });
        data.localStream = null;
      }
    };

    const attachMediaStreamToLocalVideo = () => {
      data.localVideo.srcObject = data.localStream;
      data.localVideo.muted = true;
      data.localVideo.style.display = "inline-block";
    };

    const attachMediaStreamToPeerConnection = () => {
      data.localStream.getTracks().forEach((track) => {
        data.rtcPeerConn.addTrack(track);
      });
    };

    const createOffer = () => {
      var offerOpts = {
        offerToReceiveAudio: true,
        offerToReceiveVideo: true,
      };
      data.rtcPeerConn
        .createOffer(offerOpts)
        .then((desc) => {
          data.rtcPeerConn.setLocalDescription(desc);
          var proxySdpMessage = {
            messageId: "PROXY",
            type: "sdp",
            fromId: props.ownerId,
            toId: data.contactId,
            messageData: {
              sdp: desc,
            },
          };
          data.wsConn.send(JSON.stringify(proxySdpMessage));
        })
        .catch((err) => {
          console.log(
            `createOffer failed, error name: ${err.name}, error message: ${err.message}`
          );
        });
    };

    const createAnswer = () => {
      data.rtcPeerConn
        .createAnswer()
        .then((desc) => {
          data.rtcPeerConn.setLocalDescription(desc);
          var proxySdpMessage = {
            messageId: "PROXY",
            type: "sdp",
            fromId: props.ownerId,
            toId: data.contactId,
            messageData: {
              sdp: desc,
            },
          };
          data.wsConn.send(JSON.stringify(proxySdpMessage));
        })
        .catch((err) => {
          console.log(
            `createAnswer failed, error name: ${err.name}, error message: ${err.message}`
          );
        });
    };

    const startCall = async (isInitiator, contactId) => {
      contactId = data.contactId;
      createRtcPeerConnection();
      data.localStream = await getLocalMediaStream();
      attachMediaStreamToLocalVideo();
      attachMediaStreamToPeerConnection();
      if (isInitiator) {
        var startCallMessage = {
          messageId: "PROXY",
          type: "start_call",
          fromId: props.ownerId,
          toId: contactId,
        };
        data.wsConn.send(JSON.stringify(startCallMessage));
      } else {
        var receiveCallMessage = {
          messageId: "PROXY",
          type: "receive_call",
          fromId: props.ownerId,
          toId: contactId,
        };
        data.wsConn.send(JSON.stringify(receiveCallMessage));
      }
    };

    const endCall = () => {
      closeLocalMediaStream();
      closeRtcPeerConnection();
      data.localVideo.style.display = "none";
      data.remoteVideo.style.display = "none";
      data.remoteStream = null;
      contactId = null;
    };

    const handleOfferSdp = (val) => {
      data.rtcPeerConn
        .setRemoteDescription(new RTCSessionDescription(val))
        .then(() => {
          createAnswer();
        })
        .catch((err) => {
          console.log("rtcPeerConn setRemoteDescription failed", err);
        });
    };

    const handleAnswerSdp = (val) => {
      data.rtcPeerConn.setRemoteDescription(new RTCSessionDescription(val));
    };

    const handleCandidate = (val) => {
      data.rtcPeerConn.addIceCandidate(new RTCIceCandidate(val));
    };

    const login = () => {
      if (props.ownerId === "") {
        ElMessage.warning("用户名不能为空");
        return;
      }
      data.wsConn.onopen = () => {
        console.log("连接信令服务器成功");
      };
      data.wsConn.onclose = () => {
        console.log("连接信令服务器断开");
      };
      data.wsConn.onerror = (error) => {
        console.log("连接信令服务器失败，错误信息：", error);
      };
      data.wsConn.onmessage = (event) => {
        var message = JSON.parse(event.data);
        var messageObj = message.obj; // 后端message的该字段命名为obj
        if (messageObj.messageId === "CURRENT_PEERS") {
          console.log(
            "获取CURRENT_PEERS当前在线用户列表，curContactList:",
            messageObj.messageData.curContactList
          );
          data.curContactList = messageObj.messageData.curContactList;
        } else if (messageObj.messageId === "PEER_JOIN") {
          console.log(
            "接受到PEER_JOIN消息，contactId:",
            messageObj.messagecontactId
          );
          data.curContactList.push(messageObj.messagecontactId);
        } else if (messageObj.messageId === "PEER_LEAVE") {
          console.log(
            "获取PEER_LEAVE消息，contactId:",
            messageObj.messagecontactId
          );
          for (let i = 0; i < data.curContactList.length; i++) {
            if (
              data.curContactList[i].contactId === messageObj.messagecontactId
            ) {
              if (contactId === messageObj.messagecontactId) {
                endCall();
              }
            }
          }
        } else if (messageObj.messageId === "PROXY") {
          console.log("接收到PROXY消息：", message);
          if (messageObj.type === "start_call") {
            startCall(false, messageObj.contactId);
          } else if (messageObj.type === "receive_call") {
            createOffer();
          } else if (messageObj.type === "sdp") {
            if (messageObj.messageData.sdp.type === "offer") {
              handleOfferSdp(messageObj.messageData.sdp);
            } else if (messageObj.messageData.sdp.type === "answer") {
              handleAnswerSdp(messageObj.messageData.sdp);
            } else {
              console.log("不支持的sdp类型");
            }
          } else if (messageObj.type === "candidate") {
            handleCandidate(messageObj.messageData.candidate);
          } else {
            console.log("不支持的proxy类型");
          }
        }
      };
    };

    const logout = () => {
      if (data.wsConn) {
        data.wsConn.close();
        data.wsConn = null;
      }
      endCall();
    };

    return {
      // prepareAVBase,
      createRtcPeerConnection,
      closeRtcPeerConnection,
      getLocalMediaStream,
      closeLocalMediaStream,
      attachMediaStreamToLocalVideo,
      attachMediaStreamToPeerConnection,
      createOffer,
      createAnswer,
      startCall,
      endCall,
      handleOfferSdp,
      handleAnswerSdp,
      handleCandidate,
      login,
      logout,
      emitSession,
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
  height: 360px;
  width: 700px;
}

.modal-footer {
  height: 80px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-footer-btn {
  background-color: rgb(252, 210.9, 210.9);
}
</style>
