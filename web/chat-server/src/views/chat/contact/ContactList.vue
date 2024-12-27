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
                <button class="icon-btn">
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
                <button class="icon-btn" @click="handleToOwnInfo">
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
                                  <el-dropdown-item @click="handleAgree(newContact.contact_id)">同意</el-dropdown-item>
                                  <el-dropdown-item @click="handleReject(newContact.contact_id)"> 拒绝 </el-dropdown-item>
                                  <el-dropdown-item @click="handleBlack(newContact.contact_id)"> 拉黑 </el-dropdown-item>
                                </el-dropdown-menu>
                              </template>
                            </el-dropdown>
                          </li>
                        </ul>
                      </el-scrollbar>
                    </div>
                  </template>
                  <template v-slot:footer>
                    <div class="newcontact-modal-footer">
                    </div>
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
        <el-container class="chat-container">
          <el-header>
            <h2 class="chat-name"></h2>
          </el-header>
          <el-main>
            <div class="chat-screen"></div>
            <div class="tool-bar">
              <div class="tool-bar-left">
                <el-tooltip
                  effect="customized"
                  content="表情包"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733502796507"
                      class="sticker-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="1555"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M504.32 31.872a472.256 472.256 0 1 1 0 944.512 472.256 472.256 0 0 1 0-944.512z m0 63.36a408.96 408.96 0 1 0 0 817.856 408.96 408.96 0 0 0 0-817.92z m228.864 487.808v0.192a217.856 217.856 0 1 1-435.712 0V583.04h435.712zM370.496 321.536a73.024 73.024 0 1 1 0 146.048 73.024 73.024 0 0 1 0-146.048z m289.664 0a73.024 73.024 0 1 1 0 146.048 73.024 73.024 0 0 1 0-146.048z"
                        fill="#2c2c2c"
                        p-id="1556"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="文件上传"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503065264"
                      class="upload-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="2430"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M543.7 157v534c0 16.6-13.4 30-30 30s-30-13.4-30-30V157c0-16.6 13.4-30 30-30 16.5 0 30 13.4 30 30z"
                        fill=""
                        p-id="2431"
                      ></path>
                      <path
                        d="M323.1 331c11.8 11.8 30.7 11.8 42.5 0l119.9-119.9c15.6-15.6 40.9-15.6 56.6 0L662 331c11.7 11.7 30.7 11.7 42.4 0s11.7-30.7 0-42.4L541.7 126.1c-15.6-15.6-41-15.6-56.6 0L323 288.6c-11.6 11.8-11.6 30.7 0.1 42.4zM853.7 913h-680c-33.1 0-60-26.9-60-60V583.7c0-16.4 12.8-30.2 29.2-30.7 16.9-0.4 30.8 13.2 30.8 30v240c0 16.6 13.4 30 30 30h620c16.6 0 30-13.4 30-30V583.7c0-16.4 12.8-30.2 29.2-30.7 16.9-0.4 30.8 13.2 30.8 30v270c0 33.1-26.9 60-60 60z"
                        fill=""
                        p-id="2432"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="聊天记录"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733504061769"
                      class="record-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="5492"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M695.04 194.32H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h596.96c18.32 0 33.16 14.85 33.16 33.16 0 18.31-14.84 33.16-33.16 33.16zM298.97 393.3H96.19c-17.27 0-31.27-14-31.27-31.27v-3.79c0-17.27 14-31.27 31.27-31.27h202.78c17.27 0 31.27 14 31.27 31.27v3.79c-0.01 17.28-14.01 31.27-31.27 31.27zM230.74 592.29H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h132.66c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16zM230.74 791.28H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h132.66c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16zM728.2 691.78H595.55c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16H728.2c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16z"
                        p-id="5493"
                      ></path>
                      <path
                        d="M562.38 658.62V525.96c0-18.32 14.85-33.16 33.16-33.16 18.32 0 33.16 14.85 33.16 33.16v132.66c0 18.32-14.85 33.16-33.16 33.16-18.31 0-33.16-14.85-33.16-33.16z"
                        p-id="5494"
                      ></path>
                      <path
                        d="M960.35 625.45c0 183.16-148.48 331.64-331.64 331.64S297.07 808.62 297.07 625.45s148.48-331.64 331.64-331.64 331.64 148.48 331.64 331.64zM628.71 360.14c-146.53 0-265.31 118.79-265.31 265.31s118.79 265.31 265.31 265.31 265.31-118.79 265.31-265.31-118.78-265.31-265.31-265.31z"
                        p-id="5495"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="全文复制"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503137487"
                      class="copy-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="3442"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M394.666667 106.666667h448a74.666667 74.666667 0 0 1 74.666666 74.666666v448a74.666667 74.666667 0 0 1-74.666666 74.666667H394.666667a74.666667 74.666667 0 0 1-74.666667-74.666667V181.333333a74.666667 74.666667 0 0 1 74.666667-74.666666z m0 64a10.666667 10.666667 0 0 0-10.666667 10.666666v448a10.666667 10.666667 0 0 0 10.666667 10.666667h448a10.666667 10.666667 0 0 0 10.666666-10.666667V181.333333a10.666667 10.666667 0 0 0-10.666666-10.666666H394.666667z m245.333333 597.333333a32 32 0 0 1 64 0v74.666667a74.666667 74.666667 0 0 1-74.666667 74.666666H181.333333a74.666667 74.666667 0 0 1-74.666666-74.666666V394.666667a74.666667 74.666667 0 0 1 74.666666-74.666667h74.666667a32 32 0 0 1 0 64h-74.666667a10.666667 10.666667 0 0 0-10.666666 10.666667v448a10.666667 10.666667 0 0 0 10.666666 10.666666h448a10.666667 10.666667 0 0 0 10.666667-10.666666v-74.666667z"
                        fill="#000000"
                        p-id="3443"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
              </div>
              <div class="tool-bar-right">
                <el-tooltip
                  effect="customized"
                  content="音视频通话"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503700535"
                      class="av-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="4492"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M790.207709 1023.317561c-100.48917-0.05687-302.832389-33.89448-528.321671-260.00933C-57.722981 442.903032-9.212929 154.458736 25.02277 119.995557L114.194824 30.709763c19.506387-19.563257 47.372654-30.709763 76.319449-30.709763 28.662446 0 56.073753 10.975897 75.23892 30.141064l3.980896 4.606465 131.881373 176.865489c35.145618 52.377208 33.32578 108.564701-4.720205 146.781295l-39.012773 39.069643c11.942686 71.087415 42.31123 113.398645 87.181606 158.439632l5.686993 5.686993c51.865378 52.092858 96.678885 97.076974 174.021993 103.730756l38.899033-38.955903a99.522381 99.522381 0 0 1 71.883595-30.368544c24.169721 0 49.419971 8.41675 73.020993 24.340331l178.002888 133.303121c21.212485 14.558703 34.918138 38.728424 37.477285 66.253471a113.853604 113.853604 0 0 1-33.26891 89.513274l-89.058314 89.285793c-22.179274 22.236144-85.304898 24.624681-111.465068 24.624681h-0.056869zM190.628013 88.091525a19.278907 19.278907 0 0 0-13.421304 5.402644L94.290348 176.63801c-4.549595 22.861713-44.984116 247.554815 230.607575 523.885815 202.684439 203.196268 377.50261 233.507942 463.774297 233.507942 30.652893 0 50.898589-3.753416 58.121071-5.402643l80.982784-82.006443a26.160169 26.160169 0 0 0 7.67744-18.539598l-178.457847-135.293568c-4.151505-2.786627-12.568255-7.677441-20.302566-7.677441a13.478174 13.478174 0 0 0-10.009108 3.980895l-65.969121 66.196601-18.653338-0.17061c-125.227591-1.080529-193.812729-69.950017-254.322337-130.743974l-5.686993-5.630123c-52.490947-52.661557-102.763968-117.20893-115.445963-232.199934l-2.388537-21.155614L333.826502 295.609908c8.41675-8.41675 1.990448-22.349883-4.833944-32.586471L200.750861 91.105631a17.515939 17.515939 0 0 0-10.122848-3.014106z m350.603132 312.159058c-44.131067 0-79.959125-34.235699-79.959125-76.319449V170.609797c0-42.08375 35.828057-76.376319 79.959125-76.376319h292.311452c37.136066 0 68.812618 77.968677 77.627457 111.863156 8.1324-4.606465 14.103743-8.07553 15.923581-9.269799 8.75797-5.743863 18.937687-62.670665 29.458625-62.670665a53.457736 53.457736 0 0 1 25.36399 6.426303 56.130623 56.130623 0 0 1 29.003666 49.87493v121.303566c0 21.496834-11.373986 40.775741-29.572365 50.443629a52.547817 52.547817 0 0 1-24.681551 6.141953c-10.577807 0-21.041875-56.983672-29.970454-62.955015-2.331667-1.421748-8.814839-5.118294-17.686549-10.179718-11.089637 30.368544-41.515051 105.038765-75.40953 105.038765H541.231145z m283.326003-88.944574V183.178052H550.273464v128.127957h274.283684z"
                        fill="#666666"
                        p-id="4493"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
              </div>
            </div>
          </el-main>
          <el-footer>
            <div class="chat-input">
              <el-input
                v-model="chatMessage"
                type="textarea"
                show-word-limit
                maxlength="500"
                :autosize="{ minRows: 7.9, maxRows: 7 }"
                placeholder="请输入内容"
              />
            </div>
            <div class="chat-send">
              <el-button class="send-btn">发送</el-button>
            </div>
          </el-footer>
        </el-container>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "vuex";
import Modal from "@/components/Modal.vue";
import SmallModal from "@/components/SmallModal.vue";
import axios from "axios";
import { ElMessage } from "element-plus";
export default {
  name: "ContactList",
  components: {
    Modal,
    SmallModal,
  },
  setup() {
    const data = reactive({
      chatMessage: "",
      chatName: "",
      userInfo: {
        uuid: "",
        nickname: "",
        telephone: "",
      },
      contactSearch: "",
      createGroupReq: {
        owner_id: "",
        name: "",
        notice: "",
        add_mode: null,
        avatar: "",
      },
      isCreateGroupModalVisible: false,
      isApplyContactModalVisible: false,
      isNewContactModalVisible: false,
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
        } catch (error) {
          console.error("反序列化用户信息时出错:", error);
          data.userInfo = {};
        }
      } else {
        data.userInfo = {};
      }
    });
    const router = useRouter();
    const store = useStore();
    const handleToOwnInfo = () => {
      router.push("/chat/owninfo");
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

    const handleToSessionList = () => {
      router.push("/chat/sessionList");
    };

    const handleToChatUser = async (user) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: user.user_id,
        };
        const rsp = await axios.post(store.state.backendUrl + "/session/checkOpenSessionAllowed", req);
        if (rsp.data.code == 200) {
          if (rsp.data.data == true) {
            router.push("/chat/" + user.user_id);
          } else {
            ElMessage.warning(rsp.data.message);
            console.warning(rsp.data.message);
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
        const rsp = await axios.post(store.state.backendUrl + "/contact/passContactApply", req);
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          data.newContactList = data.newContactList.filter(c => c.contact_id !== contactId);
        } else {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    return {
      ...toRefs(data),
      router,
      handleToOwnInfo,
      handleCreateGroup,
      showCreateGroupModal,
      closeCreateGroupModal,
      quitCreateGroupModal,
      showApplyContactModal,
      quitApplyContactModal,
      closeApplyContactModal,
      showNewContactModal,
      quitNewContactModal,
      handleShowUserList,
      handleHideUserList,
      handleShowMyGroupList,
      handleHideMyGroupList,
      handleShowMyJoinedGroupList,
      handleHideMyJoinedGroupList,
      handleToSessionList,
      handleToChatUser,
      handleToChatGroup,
      handleNewContactList,
      handleAgree,
    };
  },
};
</script>

<style scoped>
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
  /*background-color:aqua;*/
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
</style>