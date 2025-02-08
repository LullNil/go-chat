<template>
  <div :class="['chat-container', { 'dark-theme': isDarkTheme }]">
    <div class="main-content">
      <MainSidebar
        :isDarkTheme="isDarkTheme"
        :isSidebarOpen="isSidebarOpen"
        @toggle-theme="toggleTheme"
        @toggle-middle-panel="toggleMiddlePanel"
      />

      <MiddlePanel
        :isMiddleOpen="isMiddlePanelOpen"
        @toggle-middle-panel="toggleMiddlePanel"
      />

      <ChatBox
        :messages="messages"
        @send-message="sendMessage"
        @toggle-middle-panel="toggleMiddlePanel"
        @toggle-sidebar="toggleSidebar"
      />
    </div>
  </div>
</template>

<script>
import MainSidebar from "./MainSidebar.vue";
import MiddlePanel from "./MiddlePanel.vue";
import ChatBox from "./ChatBox.vue";

export default {
  components: {
    MainSidebar,
    MiddlePanel,
    ChatBox,
  },

  data() {
    return {
      ws: null,
      messages: [],
      isSidebarOpen: false,
      // Читаем значения из localStorage, если они там есть, иначе используем дефолтные
      isMiddlePanelOpen: localStorage.getItem("isMiddlePanelOpen")
        ? JSON.parse(localStorage.getItem("isMiddlePanelOpen"))
        : false,
      isDarkTheme: localStorage.getItem("isDarkTheme")
        ? JSON.parse(localStorage.getItem("isDarkTheme"))
        : false,
    };
  },

  mounted() {
    this.ws = new WebSocket("ws://localhost:8081/ws");

    this.ws.onmessage = (event) => {
      setTimeout(() => {
        try {
          const msg = JSON.parse(event.data);
          this.messages.push(msg);
          this.$nextTick(this.scrollToBottom);
        } catch (error) {
          console.error("Ошибка парсинга JSON:", error);
        }
      }, 1000);
    };

    this.ws.onopen = () => console.log("WebSocket подключен");
    this.ws.onerror = (error) => console.error("WebSocket ошибка:", error);
    this.ws.onclose = () => console.log("WebSocket закрыт");
  },

  methods: {
    sendMessage(msg) {
      this.messages.push(msg);
      this.ws.send(JSON.stringify(msg));
    },

    toggleTheme() {
      this.isDarkTheme = !this.isDarkTheme;
      // Set theme condition to localStorage
      localStorage.setItem("isDarkTheme", JSON.stringify(this.isDarkTheme));
    },

    toggleMiddlePanel() {
      this.isMiddlePanelOpen = !this.isMiddlePanelOpen;
      // Set middle panel condition to localStorage
      localStorage.setItem(
        "isMiddlePanelOpen",
        JSON.stringify(this.isMiddlePanelOpen)
      );
    },

    toggleSidebar() {
      this.isSidebarOpen = !this.isSidebarOpen;
    },
  },
};
</script>
