<template>
  <div class="chat-area">
    <!-- Left Side Panel -->
    <div :class="['middle-panel', { closed: !isMiddleOpen }]">
      <!-- Button to toggle panel visibility -->
      <div class="close-mid-btn">
        <span @click="toggleMiddlePanel" class="bi bi-list"></span>
      </div>

      <!-- Full chat list when panel is open -->
      <div v-if="isMiddleOpen" class="chat-list">
        <button
          @click="selectChat('general')"
          :class="{ active: activeChat === 'general' }"
        >
          <i class="bi bi-people"></i> Общий чат
        </button>
        <button
          @click="selectChat('echo')"
          :class="{ active: activeChat === 'echo' }"
        >
          <i class="bi bi-broadcast"></i> Эхо-чат
        </button>
      </div>

      <!-- Minimized chat list when panel is closed -->
      <div v-else class="chat-list-minimized">
        <button
          @click="selectChat('general')"
          :class="{ active: activeChat === 'general' }"
        >
          <i class="bi bi-people"></i>
        </button>
        <button
          @click="selectChat('echo')"
          :class="{ active: activeChat === 'echo' }"
        >
          <i class="bi bi-broadcast"></i>
        </button>
      </div>

      <!-- Bottom area for settings and theme toggle -->
      <div class="bottom-area" v-if="isMiddleOpen">
        <button class="bi bi-gear"></button>
        <button class="bi bi-person-circle"></button>
        <button
          class="theme-toggle"
          @click="toggleTheme"
          aria-label="Переключить тему"
        >
          <svg
            class="theme-icon"
            :class="{ 'animate-icon': animateIcon }"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path :d="isDarkTheme ? moonIcon : sunIcon" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Chat Box -->
    <div class="chat-box">
      <div v-if="activeChat" class="chat-messages custom-scroll">
        <div class="chat-title">
          <!-- Button for toggling panel on small screens -->
          <button @click="toggleMiddlePanel" class="menu-toggle">
            <span class="bi bi-list"></span>
          </button>
          <h2>{{ chatTitle }}</h2>
        </div>

        <!-- Message list -->
        <div ref="chatContainer" class="chat-messages custom-scroll">
          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="chat-message"
            :class="msg.user === 'Вы' ? 'sent' : 'received'"
          >
            <p class="chat-text">{{ msg.text }}</p>
          </div>
        </div>

        <!-- Input field and send button -->
        <div class="chat-input">
          <textarea
            class="input-text custom-scroll"
            v-model="newMessage"
            @keydown.enter.exact.prevent="handleSendMessage"
            @keydown.enter.ctrl.exact="addNewLine"
            placeholder="Введите сообщение..."
            rows="1"
            ref="textarea"
            @input="adjustTextareaHeight"
          ></textarea>
          <button class="btn-circle" @click="handleSendMessage">
            <svg class="send-icon" viewBox="0 0 24 24" fill="none">
              <path
                d="M5.74 15.75L3.5 20.5L20.5 12L3.5 3.5L5.74 8.25L12.5 12L5.74 15.75Z"
                fill="currentColor"
                stroke="currentColor"
                stroke-width="1.2"
                stroke-linejoin="round"
              />
            </svg>
          </button>
        </div>
      </div>
      <!-- Placeholder when no chat is selected -->
      <div v-else class="chat-placeholder">
        <p>Выберите чат для начала общения</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ChatComponent",
  props: {
    isDarkTheme: {
      type: Boolean,
      required: true,
    },
  },

  data() {
    return {
      // Chat state
      activeChat: "", // "general" or "echo", either empty string (chat not selected)
      messages: [],
      ws: null,
      newMessage: "",

      // Middle panel visibility (saves in localStorage)
      isMiddleOpen: localStorage.getItem("isMiddlePanelOpen")
        ? JSON.parse(localStorage.getItem("isMiddlePanelOpen"))
        : false,

      // Icons for theme toggle
      sunIcon: `M12 3v1
                  m0 16v1
                  m9-9h-1
                  M4 12H3
                  m15.364-6.364l-.707.707
                  M6.343 17.657l-.707.707
                  m12.728 0l-.707-.707
                  M6.343 6.343l-.707-.707
                  M12 8a4 4 0 110 8 4 4 0 010-8z`,
      moonIcon: `M21 12.79A9 9 0 1111.21 3
                   a7 7 0 009.79 9.79z`,
      animateIcon: false,
    };
  },

  mounted() {
    // Add event listeners for mouse and keyboard events
    window.addEventListener("mouseup", this.handleMouseButton);
    window.addEventListener("keydown", this.handleKeyCombination);
    // Apply theme on component mount
    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },

  beforeUnmount() {
    window.removeEventListener("mouseup", this.handleMouseButton);
    window.removeEventListener("keydown", this.handleKeyCombination);
    this.disconnectWebSocket();
  },

  computed: {
    // Returns the chat title based on the active chat type
    chatTitle() {
      switch (this.activeChat) {
        case "general":
          return "Общий чат";
        case "echo":
          return "Эхо-чат";
        default:
          return "Чат";
      }
    },
  },

  methods: {
    // Switch between Dark and Light themes
    toggleTheme() {
      this.$emit("toggle-theme");
    },

    // Toggle visibility of the middle panel
    toggleMiddlePanel() {
      this.isMiddleOpen = !this.isMiddleOpen;
      localStorage.setItem(
        "isMiddlePanelOpen",
        JSON.stringify(this.isMiddleOpen)
      );
    },

    // Select a chat mode and establish a new WebSocket connection
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.closeChat();
      this.activeChat = chatType;
      this.messages = [];
      this.connectWebSocket();
    },

    // Close the current chat: disconnect WS and clear messages
    closeChat() {
      this.disconnectWebSocket();
      this.activeChat = "";
      this.messages = [];
    },

    // Establish a WebSocket connection based on the active chat
    connectWebSocket() {
      let url = "";
      if (this.activeChat === "general") {
        url = "ws://localhost:8081/ws/general";
      } else if (this.activeChat === "echo") {
        url = "ws://localhost:8081/ws/echo";
      }
      if (url) {
        this.ws = new WebSocket(url);

        // WebSocket open event handler
        this.ws.onopen = () => {
          console.log(`WebSocket connected: ${this.activeChat}`);
        };

        // WebSocket message event: parse and add message to the list
        this.ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data);
            if (this.activeChat === "echo") {
              // For echo chat, add a delay before displaying the received message
              setTimeout(() => {
                this.messages.push(msg);
                // Auto-scroll to the new message after it is rendered
                this.$nextTick(() => {
                  this.scrollToBottom();
                });
              }, 1000);
            } else {
              this.messages.push(msg);
              // Auto-scroll for received message
              this.$nextTick(() => {
                this.scrollToBottom();
              });
            }
          } catch (error) {
            console.error("Error parsing JSON:", error);
          }
        };

        this.ws.onerror = (error) => console.error("WebSocket error:", error);
        this.ws.onclose = () => console.log("WebSocket closed");
      }
    },

    // Disconnect the current WebSocket connection if it exists
    disconnectWebSocket() {
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
    },

    // Send a message via WebSocket and update the UI accordingly
    handleSendMessage() {
      if (!this.newMessage.trim()) return;
      const msg = { user: "Вы", text: this.newMessage };

      if (this.activeChat === "echo") {
        // Immediately display the sent message in echo chat
        this.messages.push(msg);
      }
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(msg));
      } else {
        console.error("WebSocket не открыт.");
      }
      this.newMessage = "";
      this.$nextTick(() => {
        this.scrollToBottom();
        this.adjustTextareaHeight();
      });
    },

    // Add a new line in the textarea on Ctrl+Enter
    addNewLine(event) {
      event.preventDefault();
      this.newMessage += "\n";
      this.$nextTick(() => {
        this.adjustTextareaHeight();
        const textarea = this.$refs.textarea;
        if (textarea) {
          textarea.scrollTop = textarea.scrollHeight;
        }
      });
    },

    // Dynamically adjust the textarea height based on its content
    adjustTextareaHeight() {
      const textarea = this.$refs.textarea;
      if (textarea) {
        textarea.style.height = "auto";
        textarea.style.height = `${Math.min(textarea.scrollHeight, 150)}px`;
      }
    },

    // Scroll the chat container to the bottom
    scrollToBottom() {
      const chatContainer = this.$refs.chatContainer;
      if (chatContainer) {
        chatContainer.scrollTop = chatContainer.scrollHeight;
      }
    },

    // Handle mouse button events (e.g., кнопка "Назад" для закрытия чата)
    handleMouseButton(e) {
      if (e.button === 3) {
        e.preventDefault();
        this.closeChat();
      }
    },

    // Handle key combination (Ctrl+Shift+X) to close the chat
    handleKeyCombination(e) {
      if (e.ctrlKey && e.shiftKey && e.code === "KeyX") {
        e.preventDefault();
        this.closeChat();
      }
    },
  },
};
</script>

<style>
/* Import original styles to preserve layout and appearance */
@import "@/styles/chat.css";
@import "@/styles/midpanel.css";

/* Additional styling for combined component */
.chat-area {
  display: flex;
  flex: 1;
  gap: 8px;
  height: 100%;
  overflow-y: auto;
  /* scroll-behavior: smooth; */
}
</style>
