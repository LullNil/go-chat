<template>
  <div :class="['sidebar', { 'sidebar-open': isSidebarOpen }]">
    <button @click="$emit('toggle-theme')" class="btn-theme">
      <span :class="isDarkTheme ? 'bi bi-sun' : 'bi bi-moon'"></span>
    </button>

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
    isSidebarOpen: Boolean,
  },
};
</script>

<style scoped>
@import "@/styles/sidebar.css";
</style>
