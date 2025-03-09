<template>
  <div class="chat-area">
    <!-- Display ChatNavigation on desktop or when no chat is active on mobile -->
    <ChatNavigation
      v-if="!isMobile || (isMobile && !activeChat)"
      :activeChat="activeChat"
      :isDarkTheme="isDarkTheme"
      :isMobile="isMobile"
      @chat-selected="selectChat"
      @toggle-theme="toggleTheme"
      ref="navigation"
    />

    <!-- Display ChatWindow on desktop or when a chat is active on mobile -->
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
    // Determine if the device is mobile based on window width
    isMobile() {
      return this.windowWidth < 768;
    },
  },
  watch: {
    // Update background when theme changes
    isDarkTheme(newVal) {
      this.updatePageBackground(newVal);
    },
  },
  mounted() {
    // Set initial page background and add event listeners for resizing and chat closing
    this.updatePageBackground(this.isDarkTheme);
    window.addEventListener("resize", this.handleResize);
    window.addEventListener("keydown", this.handleCloseChatWindow);
    window.addEventListener("mousedown", this.handleCloseChatWindow);
  },
  beforeUnmount() {
    // Remove event listeners on component unmount
    window.removeEventListener("resize", this.handleResize);
    window.removeEventListener("keydown", this.handleCloseChatWindow);
    window.removeEventListener("mousedown", this.handleCloseChatWindow);
  },
  methods: {
    // Update the page background and meta theme-color based on the current theme
    updatePageBackground(isDark) {
      const bgColor = isDark ? "#232526" : "#ffffff";
      document.body.style.background = bgColor;
      const themeMeta = document.querySelector('meta[name="theme-color"]');
      if (themeMeta) {
        themeMeta.setAttribute("content", bgColor);
      }
    },

    // Update windowWidth on resize event
    handleResize() {
      this.windowWidth = window.innerWidth;
    },

    // Set the active chat type when a chat is selected
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.activeChat = chatType;
    },

    // Emit event to toggle the theme
    toggleTheme() {
      this.$emit("toggle-theme");
    },

    // Close the active chat
    closeChat() {
      this.activeChat = "";
    },

    // Handle events to close the chat window:
    // Ctrl+Shift+X on keyboard or middle mouse button click
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
