<template>
  <div :class="['middle-panel', { closed: !isMiddleOpen }]">
    <span @click="$emit('toggle-middle-panel')" class="bi bi-list"></span>
  </div>
</template>

<script>
export default {
  props: {
    isMiddleOpen: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped>
@import "@/styles/panel.css";
</style>
