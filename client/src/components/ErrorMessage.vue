<template>
  <transition name="fade-error">
    <div v-if="message" class="error-message">
      {{ message }}
    </div>
  </transition>
</template>

<script>
export default {
  name: "ErrorMessage",
  props: {
    message: {
      type: String,
      default: "",
    },
  },
};
</script>

<style scoped>
@import "@/styles/authpage.css";
</style>
