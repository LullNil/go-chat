<template>
  <div :class="['chat-container', { 'dark-theme': isDarkTheme }]">
    <div class="main-content">
      <MainSidebar
        :isDarkTheme="isDarkTheme"
        @toggle-theme="toggleTheme"
        @toggle-middle-panel="toggleMiddlePanel"
      />
      <MiddlePanel
        :isMiddleOpen="isMiddlePanelOpen"
        :activeChat="activeChat"
        @toggle-middle-panel="toggleMiddlePanel"
        @select-chat="selectChat"
        @close-chat="closeChat"
      />
      <ChatBox
        :messages="activeMessages"
        :activeChat="activeChat"
        @send-message="handleSendMessage"
        @toggle-middle-panel="toggleMiddlePanel"
      />
    </div>
  </div>
</template>

<script>
import { ref, watch, onBeforeUnmount } from "vue";
import MainSidebar from "./MainSidebar.vue";
import MiddlePanel from "./MiddlePanel.vue";
import ChatBox from "./ChatBox.vue";
import useWebSocket from "@/composables/useWebSocket";

export default {
  components: { MainSidebar, MiddlePanel, ChatBox },

  setup() {
    const { isDarkTheme, toggleTheme } = useTheme();
    const { isMiddlePanelOpen, toggleMiddlePanel } = useMiddlePanel();
    const {
      activeChat,
      activeMessages,
      selectChat,
      closeChat,
      handleSendMessage,
    } = useChat();

    return {
      isDarkTheme,
      toggleTheme,
      isMiddlePanelOpen,
      toggleMiddlePanel,
      activeChat,
      activeMessages,
      selectChat,
      closeChat,
      handleSendMessage,
    };
  },
};

// Composables
function useTheme() {
  const isDarkTheme = ref(
    JSON.parse(localStorage.getItem("isDarkTheme") || false)
  );

  const toggleTheme = () => {
    isDarkTheme.value = !isDarkTheme.value;
    localStorage.setItem("isDarkTheme", JSON.stringify(isDarkTheme.value));
  };

  return { isDarkTheme, toggleTheme };
}

function useMiddlePanel() {
  const isMiddlePanelOpen = ref(
    JSON.parse(localStorage.getItem("isMiddlePanelOpen") || false)
  );

  const toggleMiddlePanel = () => {
    isMiddlePanelOpen.value = !isMiddlePanelOpen.value;
    localStorage.setItem(
      "isMiddlePanelOpen",
      JSON.stringify(isMiddlePanelOpen.value)
    );
  };

  return { isMiddlePanelOpen, toggleMiddlePanel };
}

function useChat() {
  const activeChat = ref("");
  const activeMessages = ref([]);
  const { connect, disconnect, send } = useWebSocket();

  watch(activeChat, (newChat) => {
    if (newChat) {
      connect(`ws://localhost:8081/ws/${newChat}`, handleMessage);
    } else {
      disconnect();
    }
  });

  const handleMessage = (msg) => {
    if (activeChat.value === "echo") {
      setTimeout(() => activeMessages.value.push(msg), 1000);
    } else {
      activeMessages.value.push(msg);
    }
  };

  const selectChat = (chatType) => {
    if (activeChat.value === chatType) return;
    activeChat.value = chatType;
    activeMessages.value = [];
  };

  const closeChat = () => {
    activeChat.value = "";
    activeMessages.value = [];
  };

  const handleSendMessage = (msg) => {
    if (activeChat.value === "echo") {
      activeMessages.value.push(msg);
    }
    send(msg);
  };

  onBeforeUnmount(() => disconnect());

  return {
    activeChat,
    activeMessages,
    selectChat,
    closeChat,
    handleSendMessage,
  };
}
</script>
