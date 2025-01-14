<template>
  <div v-if="isVisible">deleteUserModal</div>
</template>

<script>
export default {
  name: "DeleteUserModal",
  props: {
    isVisible: false,
  },
};
</script>

<style scoped>
</style>