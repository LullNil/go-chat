<template>
  <!-- Custom title bar -->
  <CustomTitlebar />

  <!-- Main chat interface -->
  <router-view />
</template>

<script>
import CustomTitlebar from "./components/CustomTitlebar.vue";

export default {
  name: "App",
  components: {
    CustomTitlebar,
  },
};
</script>
