<template>
  <div class="chat-area">
    <!-- Left Side Panel -->
    <div :class="['middle-panel', { closed: !isMiddleOpen }]">
      <div v-if="isMiddleOpen" class="middle-panel-content">
        <!-- Header section with title -->
        <div class="chat-header" ref="chatHeader">
          <div class="header-action" @click="toggleMiddlePanel">
            <span class="close-mid-btn bi bi-list"></span>
            <h3>Чаты</h3>
          </div>
        </div>

        <!-- Scrollable chat list area -->
        <div class="chat-list-scroll custom-scroll" ref="chatListScroll">
          <div class="chat-list">
            <button
              @click="selectChat('general')"
              :class="{ active: activeChat === 'general' }"
            >
              <i class="bi bi-people"></i> Общий чат
            </button>
            <!-- Example buttons (can be replaced with v-for loop) -->
            <button
              v-for="n in 10"
              :key="n"
              @click="selectChat('echo')"
              :class="{ active: activeChat === 'echo' }"
            >
              <i class="bi bi-broadcast"></i> Эхо-чат {{ n }}
            </button>
          </div>
        </div>

        <!-- Search block -->
        <transition name="slide-up">
          <div v-if="isSearchVisible" class="chat-search-toggle">
            <input
              ref="searchInput"
              type="text"
              v-model="searchQuery"
              placeholder="Поиск по чатам..."
            />
          </div>
        </transition>

        <!-- Bottom panel with action buttons -->
        <div class="bottom-area">
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
          <!-- Button to toggle search visibility -->
          <button
            @click="toggleSearch"
            aria-label="Поиск"
            class="bi bi-search"
          ></button>
        </div>
      </div>

      <!-- Minimized panel version -->
      <div v-else class="chat-list-minimized">
        <!-- Header for collapsing/expanding, always visible -->
        <div
          class="minimized-header custom-scroll"
          ref="minimizedHeader"
          @click="toggleMiddlePanel"
        >
          <span class="close-mid-btn bi bi-list"></span>
        </div>
        <!-- Container with chat selection buttons (icons centered with scrolling enabled) -->
        <div class="minimized-buttons custom-scroll" ref="minimizedButtons">
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
          <button
            v-for="v in 15"
            :key="v"
            @click="selectChat('echo')"
            :class="{ active: activeChat === 'echo' }"
          >
            <i class="bi bi-broadcast"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Chat Box -->
    <div class="chat-box">
      <div
        v-if="activeChat"
        class="chat-messages custom-scroll"
        ref="chatMessages"
      >
        <div class="chat-title">
          <!-- Button for toggling panel on small screens -->
          <button @click="toggleMiddlePanel" class="menu-toggle">
            <span class="bi bi-list"></span>
          </button>
          <h2>{{ chatTitle }}</h2>
        </div>

        <!-- Messages list -->
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
      activeChat: "",
      messages: [],
      ws: null,
      newMessage: "",
      isMiddleOpen: localStorage.getItem("isMiddlePanelOpen")
        ? JSON.parse(localStorage.getItem("isMiddlePanelOpen"))
        : false,
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
      isSearchVisible: false,
      searchQuery: "",
    };
  },
  mounted() {
    // Add global event listeners
    window.addEventListener("mouseup", this.handleMouseButton);
    window.addEventListener("keydown", this.handleKeyCombination);

    // Add scroll event listeners if refs exist
    if (this.$refs.chatListScroll) {
      this.$refs.chatListScroll.addEventListener(
        "scroll",
        this.handleChatListScroll
      );
    }
    if (this.$refs.minimizedButtons) {
      this.$refs.minimizedButtons.addEventListener(
        "scroll",
        this.handleMinimizedScroll
      );
    }

    // Initialize scroll handlers to set shadow correctly
    this.handleChatListScroll();
    if (this.$refs.minimizedButtons) {
      this.handleMinimizedScroll();
    }

    // Toggle dark theme based on prop
    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },
  beforeUnmount() {
    // Remove global event listeners
    window.removeEventListener("mouseup", this.handleMouseButton);
    window.removeEventListener("keydown", this.handleKeyCombination);

    // Remove scroll event listeners if refs exist
    if (this.$refs.chatListScroll) {
      this.$refs.chatListScroll.removeEventListener(
        "scroll",
        this.handleChatListScroll
      );
    }
    if (this.$refs.minimizedButtons) {
      this.$refs.minimizedButtons.removeEventListener(
        "scroll",
        this.handleMinimizedScroll
      );
    }

    this.disconnectWebSocket();
  },
  computed: {
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
    // Emit event to toggle theme
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

      this.$nextTick(() => {
        // Reinitialize scroll event listeners after toggling panel
        if (this.$refs.chatListScroll) {
          this.$refs.chatListScroll.addEventListener(
            "scroll",
            this.handleChatListScroll
          );
          this.handleChatListScroll();
        }
        if (this.$refs.minimizedButtons) {
          this.$refs.minimizedButtons.addEventListener(
            "scroll",
            this.handleMinimizedScroll
          );
          this.handleMinimizedScroll();
        }
      });
    },
    // Toggle search input visibility
    toggleSearch() {
      this.isSearchVisible = !this.isSearchVisible;
      if (this.isSearchVisible) {
        this.$nextTick(() => {
          this.$refs.searchInput.focus();
        });
      }
    },
    // Select a chat and initialize it
    selectChat(chatType) {
      if (this.activeChat === chatType) return;
      this.closeChat();
      this.activeChat = chatType;
      this.messages = [];
      this.connectWebSocket();
    },
    // Close the current chat
    closeChat() {
      this.disconnectWebSocket();
      this.activeChat = "";
      this.messages = [];
    },
    // Connect to the WebSocket for the selected chat
    connectWebSocket() {
      let url = "";
      if (this.activeChat === "general") {
        url = "ws://localhost:8081/ws/general";
      } else if (this.activeChat === "echo") {
        url = "ws://localhost:8081/ws/echo";
      }
      if (url) {
        this.ws = new WebSocket(url);
        this.ws.onopen = () => {
          console.log(`WebSocket connected: ${this.activeChat}`);
        };
        this.ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data);
            if (this.activeChat === "echo") {
              setTimeout(() => {
                this.messages.push(msg);
                this.$nextTick(() => {
                  this.scrollToBottom();
                });
              }, 1000);
            } else {
              this.messages.push(msg);
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
    // Disconnect from the WebSocket
    disconnectWebSocket() {
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
    },
    // Handle sending a message
    handleSendMessage() {
      if (!this.newMessage.trim()) return;
      const msg = { user: "Вы", text: this.newMessage };
      if (this.activeChat === "echo") {
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
    // Add a new line in the textarea
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
    // Adjust the height of the textarea based on content
    adjustTextareaHeight() {
      const textarea = this.$refs.textarea;
      if (textarea) {
        textarea.style.height = "auto";
        textarea.style.height = `${Math.min(textarea.scrollHeight, 150)}px`;
      }
    },
    // Scroll chat container to the bottom
    scrollToBottom() {
      const chatContainer = this.$refs.chatContainer;
      if (chatContainer) {
        chatContainer.scrollTop = chatContainer.scrollHeight;
      }
    },
    // Handle scroll event for chat list to toggle header shadow
    handleChatListScroll() {
      if (!this.$refs.chatListScroll) return;
      const scrollTop = this.$refs.chatListScroll.scrollTop;
      this.$refs.chatHeader?.classList.toggle("scrolled", scrollTop > 0);
    },
    // Handle scroll event for minimized buttons to toggle header shadow
    handleMinimizedScroll() {
      if (!this.$refs.minimizedButtons) return;
      const scrollTop = this.$refs.minimizedButtons.scrollTop;
      this.$refs.minimizedHeader?.classList.toggle("scrolled", scrollTop > 0);
    },
    // Close chat on mouse button event (e.g., middle click)
    handleMouseButton(e) {
      if (e.button === 3) {
        e.preventDefault();
        this.closeChat();
      }
    },
    // Close chat on specific key combination (Ctrl+Shift+X)
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
@import "@/styles/chat.css";
@import "@/styles/midpanel.css";

.chat-area {
  display: flex;
  flex: 1;
  gap: 8px;
  height: 100%;
  overflow-y: auto;
  overflow: hidden;
}
</style>
