<template>
  <!-- Custom title bar -->
  <CustomTitlebar />

  <!-- Main container -->
  <MainPage />
</template>

<script>
import CustomTitlebar from "./components/CustomTitlebar.vue";
import MainPage from "./pages/MainPage.vue";

export default {
  name: "App",
  components: {
    MainPage,
    CustomTitlebar,
  },
};
</script>
