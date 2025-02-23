<template>
  <!-- Auth Form Overlay -->

  <AuthForm
    v-if="!isAuthenticated"
    :isDarkTheme="isDarkTheme"
    @toggle-theme="toggleTheme"
    @login="handleLogin"
  />

  <!-- Main Chat Interface -->
  <div
    v-show="isAuthenticated"
    :class="['chat-container', { 'dark-theme': isDarkTheme }]"
  >
    <!-- Main content container -->
    <div class="main-content">
      <!-- MiddlePanel transmits the selected chat mode -->
      <MiddlePanel
        :isMiddleOpen="isMiddlePanelOpen"
        :activeChat="activeChat"
        :isDarkTheme="isDarkTheme"
        @toggle-theme="toggleTheme"
        @toggle-middle-panel="toggleMiddlePanel"
        @select-chat="selectChat"
        @close-chat="closeChat"
      />

      <!-- ChatBox receives list of messages and active chat -->
      <ChatBox
        :messages="activeMessages"
        :activeChat="activeChat"
        @send-message="sendMessage"
        @toggle-middle-panel="toggleMiddlePanel"
        @close-chat="closeChat"
      />

      <RightSidePanel />
    </div>
  </div>
</template>

<script>
import MiddlePanel from "./MiddlePanel.vue";
import ChatBox from "./ChatBox.vue";
import RightSidePanel from "./RightSidePanel.vue";
import AuthForm from "./AuthForm.vue";

export default {
  components: {
    MiddlePanel,
    ChatBox,
    RightSidePanel,
    AuthForm,
  },

  data() {
    return {
      isAuthenticated:
        Boolean(localStorage.getItem("isAuthenticated")) || false,
      showAuthForm: true, // Форма авторизации показывается при каждой загрузке страницы
      activeChat: "", // "general" или "echo", или пустая строка (чат не выбран)
      activeMessages: [],
      ws: null,
      isMiddlePanelOpen: localStorage.getItem("isMiddlePanelOpen")
        ? JSON.parse(localStorage.getItem("isMiddlePanelOpen"))
        : false,
      isDarkTheme: localStorage.getItem("isDarkTheme")
        ? JSON.parse(localStorage.getItem("isDarkTheme"))
        : false,
    };
  },

  mounted() {
    console.log("Initial auth state:", this.isAuthenticated);
    window.addEventListener("storage", (e) => {
      console.log("Storage event:", e.key, e.newValue);
    });

    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },

  methods: {
    handleLogin() {
      this.isAuthenticated = true;
      localStorage.setItem("isAuthenticated", "true");
    },

    // Method toggleTheme switches between light and dark themes
    toggleTheme() {
      this.isDarkTheme = !this.isDarkTheme;
      localStorage.setItem("isDarkTheme", JSON.stringify(this.isDarkTheme));

      document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
    },

    // Method toggleMiddlePanel toggles the visibility of the middle panel
    toggleMiddlePanel() {
      this.isMiddlePanelOpen = !this.isMiddlePanelOpen;
      localStorage.setItem(
        "isMiddlePanelOpen",
        JSON.stringify(this.isMiddlePanelOpen)
      );
    },

    // Выбор режима чата (general или echo)
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.disconnectWebSocket();
      this.activeChat = chatType;
      this.activeMessages = [];
      this.connectWebSocket();
    },

    // Закрытие чата – сбрасываем активный режим и сообщения
    closeChat() {
      this.disconnectWebSocket();
      this.activeChat = "";
      this.activeMessages = [];
    },

    /**
     * Establishes a WebSocket connection to the server
     * based on the currently selected chat type (general or echo).
     */
    connectWebSocket() {
      let url = "";
      if (this.activeChat === "general") {
        url = "ws://localhost:8081/ws/general";
      } else if (this.activeChat === "echo") {
        url = "ws://localhost:8081/ws/echo";
      }
      if (url) {
        this.ws = new WebSocket(url);

        // Event handlers
        this.ws.onopen = () =>
          console.log(`WebSocket connected: ${this.activeChat}`);
        this.ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data);
            if (this.activeChat === "echo") {
              // In echo chat, add a 1 second delay to receive the echo response
              setTimeout(() => {
                this.activeMessages.push(msg);
              }, 1000);
            } else {
              this.activeMessages.push(msg);
            }
          } catch (error) {
            console.error("Error parsing JSON:", error);
          }
        };
        this.ws.onerror = (error) => console.error("WebSocket error:", error);
        this.ws.onclose = () => console.log("WebSocket closed");
      }
    },

    disconnectWebSocket() {
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
    },

    // Отправка сообщения через WS. Для echo-чата – также добавляем сообщение сразу в список.
    sendMessage(msg) {
      if (this.activeChat === "echo") {
        // Показываем отправленное сообщение сразу в эхо-чате
        this.activeMessages.push(msg);
      }
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(msg));
      } else {
        console.error("WebSocket не открыт.");
      }
    },
  },

  beforeUnmount() {
    this.disconnectWebSocket();
  },

  handleLogin() {
    this.showAuthForm = false;
  },
};
</script>

<style scoped>
/* Анимация появления */
.auth-fade-enter-active,
.auth-fade-leave-active {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.auth-fade-enter-from,
.auth-fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
