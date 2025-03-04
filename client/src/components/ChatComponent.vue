<template>
  <div class="chat-area">
    <!-- Side Panel -->
    <div
      ref="middlePanel"
      class="middle-panel"
      :class="{ closed: !isMiddleOpen, resizing: isResizing }"
      :style="{ width: isMiddleOpen ? panelWidth + 'px' : '80px' }"
    >
      <!-- Resize Handle: Right edge of the panel -->
      <div
        class="resize-handle"
        :class="{ active: isResizing }"
        @mousedown="initResize"
      ></div>

      <!-- Expanded panel content -->
      <div v-if="isMiddleOpen" class="middle-panel-content">
        <!-- Header section -->
        <div class="chat-header" ref="chatHeader">
          <div class="header-action" @click="toggleMiddlePanel">
            <span class="close-mid-btn bi bi-list"></span>
            <h3 :class="{ 'search-active': isSearchActive }">Чаты</h3>
          </div>

          <!-- Search Container -->
          <div class="search-container" :class="{ active: isSearchActive }">
            <input
              ref="searchInput"
              type="text"
              class="search-input"
              placeholder="Поиск"
              v-model="searchQuery"
            />
            <div
              class="search-icon-wrapper"
              @click="toggleSearch"
              :class="{ active: isSearchActive }"
            >
              <svg class="search-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M15.5 14H14.71L14.43 13.73C15.41 12.59 16 11.11 16 9.5C16 5.91 13.09 3 9.5 3C5.91 3 3 5.91 3 9.5C3 13.09 5.91 16 9.5 16C11.11 16 12.59 15.41 13.73 14.43L14 14.71V15.5L19 20.49L20.49 19L15.5 14ZM9.5 14C7.01 14 5 11.99 5 9.5C5 7.01 7.01 5 9.5 5C11.99 5 14 7.01 14 9.5C14 11.99 11.99 14 9.5 14Z"
                />
              </svg>
            </div>
          </div>

          <div
            class="plus-icon-wrapper"
            :class="{ active: isMenuOpen }"
            @click="toggleMenu"
          >
            <svg class="plus-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
            </svg>

            <!-- Dropdown Menu -->
            <div class="dropdown-menu-mid" :class="{ active: isMenuOpen }">
              <div
                class="menu-item"
                v-for="(item, index) in menuItems"
                :key="index"
                @click.stop="handleMenuItemClick(item)"
              >
                <i :class="item.icon"></i>
                <span>{{ item.label }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="categories-container" ref="categoriesContainer">
          <div class="categories-scroll" @wheel="handleCategoriesScroll">
            <button
              class="category-btn"
              v-for="(category, index) in categories"
              :key="index"
              :class="{ active: activeCategory === index }"
              @mousedown.prevent
              @click="selectCategory(index)"
            >
              {{ category }}
            </button>
          </div>
        </div>

        <!-- Chat list scroll area -->
        <div class="chat-list-scroll custom-scroll" ref="chatListScroll">
          <div class="chat-list">
            <button
              @click="selectChat('general')"
              :class="{ active: activeChat === 'general' }"
            >
              <i class="bi bi-people"></i> Общий чат
            </button>
            <!-- Example chat buttons -->
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
        </div>
      </div>

      <!-- Minimized panel version -->
      <div v-else class="chat-list-minimized">
        <div
          class="minimized-header custom-scroll"
          ref="minimizedHeader"
          @click="toggleMiddlePanel"
        >
          <span class="close-mid-btn bi bi-list"></span>
        </div>
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
          <!-- Button to toggle panel on small screens -->
          <button @click="toggleMiddlePanel" class="menu-toggle">
            <span class="bi bi-list"></span>
          </button>
          <h2>{{ chatTitle }}</h2>
        </div>

        <!-- Message list -->
        <div ref="chatContainer" class="chat-messages chat-scroll">
          <div class="messages-inner">
            <div
              v-for="(msg, index) in messages"
              :key="index"
              class="chat-message"
              :class="msg.user === 'Вы' ? 'sent' : 'received'"
            >
              <p class="chat-text">{{ msg.text }}</p>
            </div>
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
      <!-- Placeholder if no chat is selected -->
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
      // Panel width is stored in localStorage and determines expansion.
      panelWidth: localStorage.getItem("panelWidth")
        ? Number(localStorage.getItem("panelWidth"))
        : 280,
      // Last expanded width to restore when expanding from minimized state
      lastExpandedWidth: 280,
      // Flag to track the resizing process
      isResizing: false,
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
      isSearchActive: false,
      searchQuery: "",
      isMenuOpen: false,
      menuItems: [
        {
          icon: "bi bi-person-plus",
          label: "Профиль",
          key: "profile",
        },
        {
          icon: "bi bi-palette",
          label: "Тема",
          key: "theme",
          action: this.toggleTheme,
        },
        {
          icon: "bi bi-gear",
          label: "Настройки",
          key: "settings",
        },
      ],
      activeCategory: 0,
      categories: ["Все", "Личные", "Группы", "Архив", "Избранное"],
      // Define thresholds
      SNAP_THRESHOLD: 240,
      MIN_THRESHOLD: 120,
    };
  },

  computed: {
    isMiddleOpen() {
      return this.panelWidth > this.MIN_THRESHOLD;
    },
    // Returns chat title based on active chat
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

  mounted() {
    // Global event listeners
    window.addEventListener("keydown", this.handleKeyCombination);
    window.addEventListener("mousedown", this.handleMouseButton); // Для закрытия чата кнопкой мыши 3
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
    this.handleChatListScroll();
    if (this.$refs.minimizedButtons) {
      this.handleMinimizedScroll();
    }
    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },

  beforeUnmount() {
    // Remove global event listeners
    window.removeEventListener("mousemove", this.onResize);
    window.removeEventListener("mouseup", this.stopResize);
    window.removeEventListener("keydown", this.handleKeyCombination);
    window.removeEventListener("mousedown", this.handleMouseButton);
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

  methods: {
    initResize(e) {
      // Prevent default behavior and start resizing
      e.preventDefault();
      this.isResizing = true;
      // Disable text selection for smoother UX
      document.body.style.userSelect = "none";
      window.addEventListener("mousemove", this.onResize);
      window.addEventListener("mouseup", this.stopResize);
    },
    onResize(e) {
      if (!this.isResizing) return;
      const panelRect = this.$refs.middlePanel.getBoundingClientRect();
      const newWidth = e.clientX - panelRect.left;
      if (newWidth < this.MIN_THRESHOLD) {
        // If dragged below MIN_THRESHOLD, set minimized width
        this.panelWidth = 80;
      } else {
        // If newWidth is between MIN_THRESHOLD and SNAP_THRESHOLD, hold the panel at SNAP_THRESHOLD;
        // otherwise, update freely.
        this.panelWidth =
          newWidth < this.SNAP_THRESHOLD ? this.SNAP_THRESHOLD : newWidth;
        this.lastExpandedWidth = this.panelWidth;
      }
    },
    stopResize() {
      if (this.isResizing) {
        this.isResizing = false;
        document.body.style.userSelect = "auto";
        window.removeEventListener("mousemove", this.onResize);
        window.removeEventListener("mouseup", this.stopResize);
        localStorage.setItem("panelWidth", this.panelWidth);

        this.$nextTick(() => {
          const handleScrollUpdate = () => {
            this.handleMinimizedScroll();
            this.handleChatListScroll();
          };

          const initScrollHandler = (ref, handler) => {
            if (ref && !ref._scrollHandlerAttached) {
              ref.addEventListener("scroll", handler);
              ref._scrollHandlerAttached = true;
            }
          };

          initScrollHandler(this.$refs.minimizedButtons, handleScrollUpdate);
          initScrollHandler(this.$refs.chatListScroll, handleScrollUpdate);

          handleScrollUpdate();
        });
      }
    },

    // Close chat on mouse button 3 (applied globally)
    handleMouseButton(e) {
      if (e.button === 3) {
        e.preventDefault();
        this.closeChat();
      }
    },

    handleMenuItemClick(item) {
      if (item.action) {
        item.action();
      }
      // Закрываем меню только для определенных действий
      if (!["theme"].includes(item.key)) {
        this.isMenuOpen = false;
      }
    },

    handleCategoriesScroll(e) {
      const container = this.$refs.categoriesContainer;
      if (container) {
        container.scrollLeft += e.deltaY;
      }
    },

    selectCategory(index) {
      this.activeCategory = index;
      // Дополнительная логика фильтрации чатов
    },

    toggleTheme() {
      this.$emit("toggle-theme");
    },
    toggleMiddlePanel() {
      // Toggle panel expansion: if expanded, save current width and minimize;
      // if minimized, restore last expanded width.
      if (this.isMiddleOpen) {
        this.lastExpandedWidth = this.panelWidth;
        this.panelWidth = 80;
      } else {
        this.panelWidth = this.lastExpandedWidth || 280;
      }
      localStorage.setItem("panelWidth", this.panelWidth);
      // Reinitialize header shadow logic after toggling panel state
      this.$nextTick(() => {
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
      this.isSearchActive = !this.isSearchActive;
      if (this.isSearchActive) {
        this.$nextTick(() => {
          this.$refs.searchInput.focus();
        });
      }
    },
    toggleMenu() {
      this.isMenuOpen = !this.isMenuOpen;
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
        console.error("WebSocket is not open.");
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

    // Method to update header shadow for expanded panel
    handleChatListScroll() {
      if (!this.$refs.chatListScroll || !this.$refs.categoriesContainer) return;
      const scrollTop = this.$refs.chatListScroll.scrollTop;
      this.$refs.categoriesContainer.classList.toggle(
        "scrolled",
        scrollTop > 0
      );
    },

    // Method to update header shadow for minimized panel
    handleMinimizedScroll() {
      if (!this.$refs.minimizedButtons) return;
      const scrollTop = this.$refs.minimizedButtons.scrollTop;
      const header = this.$refs.minimizedHeader;
      this.$refs.minimizedHeader?.classList.toggle("scrolled", scrollTop > 0);

      if (header) {
        header.classList.toggle("scrolled", scrollTop > 0);

        // Добавьте анимацию для плавного появления тени
        if (scrollTop > 0) {
          header.style.boxShadow = "0 4px 8px -2px rgba(0, 0, 0, 0.25)";
        } else {
          header.style.boxShadow = "none";
        }
      }
    },
    // Method to handle key combinations
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
</style>
