<template>
  <transition name="bounce-pop">
    <div v-if="message" class="mobile-error-message">
      <span class="error-icon">⚠</span>
      {{ message }}
    </div>
  </transition>
</template>

<script>
export default {
  name: "ErrorMessageMobile",
  props: {
    message: {
      type: String,
      default: "",
    },
  },
};
</script>
