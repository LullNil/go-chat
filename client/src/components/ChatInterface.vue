<template>
  <div class="chat-area">
    <ChatNavigation
      v-if="!isMobile || (isMobile && !activeChat)"
      :activeChat="activeChat"
      :isDarkTheme="isDarkTheme"
      :isMobile="isMobile"
      @chat-selected="selectChat"
      @toggle-theme="toggleTheme"
      ref="navigation"
    />

    <ChatWindow
      v-if="!isMobile || (isMobile && activeChat)"
      :activeChat="activeChat"
      :isMobile="isMobile"
      @toggle-sidebar="closeChat"
      @close-chat="closeChat"
      key="chatWindow"
    />
  </div>
</template>

<script>
import ChatNavigation from "./ChatNavigation.vue";
import ChatWindow from "./ChatWindow.vue";

export default {
  name: "ChatInterface",
  components: { ChatNavigation, ChatWindow },
  props: {
    isDarkTheme: { type: Boolean, required: true },
  },
  data() {
    return {
      activeChat: "",
      windowWidth: window.innerWidth,
    };
  },
  computed: {
    isMobile() {
      return this.windowWidth < 768;
    },
  },
  watch: {
    isDarkTheme(newVal) {
      this.updatePageBackground(newVal);
    },
  },

  mounted() {
    this.updatePageBackground(this.isDarkTheme);
    window.addEventListener("resize", this.handleResize);
    window.addEventListener("keydown", this.handleCloseChatWindow);
    window.addEventListener("mousedown", this.handleCloseChatWindow);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.handleResize);
    window.removeEventListener("keydown", this.handleCloseChatWindow);
    window.removeEventListener("mousedown", this.handleCloseChatWindow);
  },
  methods: {
    updatePageBackground(isDark) {
      const bgColor = isDark ? "#232526" : "#ffffff";
      document.body.style.background = bgColor;
      const themeMeta = document.querySelector('meta[name="theme-color"]');
      if (themeMeta) {
        themeMeta.setAttribute("content", bgColor);
      }
    },
    handleResize() {
      this.windowWidth = window.innerWidth;
    },
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.activeChat = chatType;
    },
    toggleTheme() {
      this.$emit("toggle-theme");
    },
    closeChat() {
      this.activeChat = "";
    },
    handleCloseChatWindow(e) {
      switch (e.type) {
        case "keydown":
          if (e.ctrlKey && e.shiftKey && e.code === "KeyX") {
            e.preventDefault();
            this.closeChat();
          }
          break;
        case "mousedown":
          if (e.button === 3) {
            e.preventDefault();
            this.closeChat();
          }
          break;
      }
    },
  },
};
</script>


