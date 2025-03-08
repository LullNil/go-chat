<template>
  <div class="chat-area">
    <!-- Chat navigation component -->
    <ChatNavigation
      :activeChat="activeChat"
      :isDarkTheme="isDarkTheme"
      @chat-selected="selectChat"
      @toggle-theme="toggleTheme"
      ref="navigation"
    />

    <!-- Chat window component -->
    <ChatWindow :activeChat="activeChat" @toggle-sidebar="toggleNavigation" />
  </div>
</template>

<script>
import ChatNavigation from "./ChatNavigation.vue";
import ChatWindow from "./ChatWindow.vue";

export default {
  name: "ChatInterface",
  components: {
    ChatNavigation,
    ChatWindow,
  },

  props: {
    isDarkTheme: {
      type: Boolean,
      required: true,
    },
  },

  data() {
    return {
      activeChat: "", // If empty - chat is closed, display placeholder
    };
  },

  mounted() {
    window.addEventListener("keydown", this.handleCloseChatWindow);
    window.addEventListener("mousedown", this.handleCloseChatWindow);
  },

  beforeUnmount() {
    window.removeEventListener("keydown", this.handleCloseChatWindow);
    window.removeEventListener("mousedown", this.handleCloseChatWindow);
  },

  methods: {
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.activeChat = chatType;
    },

    toggleTheme() {
      this.$emit("toggle-theme");
    },

    toggleNavigation() {
      if (
        this.$refs.navigation &&
        typeof this.$refs.navigation.toggleMiddlePanel === "function"
      ) {
        this.$refs.navigation.toggleMiddlePanel();
      }
    },

    // Close chat handler
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

    // Close chat method
    closeChat() {
      this.activeChat = "";
    },
  },
};
</script>
