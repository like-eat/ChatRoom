<template>
  <div class="add-friend">
    <el-input
      v-model="searchQuery"
      placeholder="请输入好友昵称进行搜索"
      clearable
      @clear="handleClear"
      @input="handleSearch"
      suffix-icon="el-icon-search"
    />
    <el-list
      v-if="filteredFriends.length > 0"
      border
      :data="filteredFriends"
      style="width: 100%">
      <el-list-item
        slot-scope="scope"
        @click="addFriend(scope.item)"
        class="friend-item">
        <el-avatar :src="scope.item.avatar" />
        <span>{{ scope.item.name }}</span>
      </el-list-item>
    </el-list>
    <el-button type="primary" @click="showNoFriendsMessage" v-if="filteredFriends.length === 0">没有找到相关好友</el-button>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';

export default {
  name: 'AddFriend',
  components: {
    ElInput,
    ElList,
    ElListItem,
    ElAvatar,
    ElButton,
  },
  setup() {
    // 模拟的好友列表
    const friends = ref([
      { name: '张三', avatar: 'https://via.placeholder.com/50' },
      { name: '李四', avatar: 'https://via.placeholder.com/50' },
      { name: '王五', avatar: 'https://via.placeholder.com/50' },
      // 更多好友...
    ]);

    // 搜索查询
    const searchQuery = ref('');

    // 根据搜索查询过滤好友列表
    const filteredFriends = computed(() => {
      if (!searchQuery.value) {
        return friends.value;
      }
      return friends.value.filter(friend => 
        friend.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      );
    });

    // 清除搜索查询
    const handleClear = () => {
      searchQuery.value = '';
    };

    // 处理搜索输入
    const handleSearch = () => {
      // 可以在这里添加搜索防抖功能
    };

    // 添加好友
    const addFriend = (friend) => {
      // 这里可以添加实际的添加好友逻辑，比如发送请求到服务器
      ElMessage.success(`成功添加好友: ${friend.name}`);
      searchQuery.value = ''; // 清空搜索框
    };

    // 显示没有找到好友的消息
    const showNoFriendsMessage = () => {
      ElMessage.info('没有找到相关好友');
    };

    return {
      searchQuery,
      filteredFriends,
      handleClear,
      handleSearch,
      addFriend,
      showNoFriendsMessage,
    };
  },
};
</script>

<style scoped>
.add-friend {
  margin: 20px;
}
.friend-item {
  display: flex;
  align-items: center;
  padding: 10px;
  cursor: pointer;
}
.friend-item .el-avatar {
  margin-right: 10px;
}
</style>