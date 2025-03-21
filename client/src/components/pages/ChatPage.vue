<template>
  <div class="chat-container">
    <!-- Display ChatSidebar on desktop or when no chat is active on mobile -->
    <ChatSidebar
      v-if="!isMobile || (isMobile && !activeChat)"
      :activeChat="activeChat"
      :isMobile="isMobile"
      @chat-selected="selectChat"
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
import ChatSidebar from "../ChatSidebar.vue";
import ChatWindow from "../ChatWindow.vue";

export default {
  name: "ChatPage",
  components: { ChatSidebar, ChatWindow },
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
  mounted() {
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
    // Set the active chat type when a chat is selected
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.activeChat = chatType;
    },
    // Close the active chat
    closeChat() {
      this.activeChat = "";
    },
    // Update windowWidth on resize event
    handleResize() {
      this.windowWidth = window.innerWidth;
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
