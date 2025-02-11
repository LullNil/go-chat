<template>
  <!-- Sidebar -->
  <div class="sidebar">
    <!-- Theme toggle button -->
    <button @click="$emit('toggle-theme')" class="btn-theme">
      <span :class="isDarkTheme ? 'bi bi-sun' : 'bi bi-moon'"></span>
    </button>

    <!-- Buttons inside the sidebar -->
    <div class="nav-icons">
      <span
        class="bi bi-chat-dots"
        @click="$emit('toggle-middle-panel')"
      ></span>
      <span class="bi bi-columns-gap"></span>
      <span class="bi bi-people"></span>
      <span class="bi bi-gear"></span>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    isDarkTheme: Boolean,
  },
};
</script>

<style scoped>
@import "@/styles/sidebar.css";
</style>
