<template>
  <!-- Middle Panel for chat navigation -->
  <div
    ref="middlePanel"
    class="middle-panel"
    :class="{ closed: !isMiddleOpen, resizing: isResizing }"
    :style="{ width: isMiddleOpen ? panelWidth + 'px' : '80px' }"
  >
    <!-- Resize handle for middle panel -->
    <div
      class="resize-handle"
      :class="{ active: isResizing }"
      @mousedown="initResize"
    ></div>

    <!-- Expanded Middle Panel Content -->
    <div v-if="isMiddleOpen" class="middle-panel-content">
      <!-- Chat header with dropdown menu trigger -->
      <div class="chat-header" ref="chatHeader">
        <div
          class="plus-icon-wrapper"
          :class="{ active: isMenuOpen }"
          @click="toggleMenu"
        >
          <svg class="plus-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
          </svg>
        </div>
        <div class="header-action">
          <h3
            @click="toggleMiddlePanel"
            :class="{ 'search-active': isSearchActive }"
          >
            Чаты
          </h3>
        </div>
        <!-- Search input field -->
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
      </div>

      <!-- Dropdown menu -->
      <div class="dropdown-menu-mid" :class="{ active: isMenuOpen }">
        <div
          v-for="(item, index) in menuItems"
          class="menu-item"
          :key="index"
          :class="{ 'logout-item': item.key === 'logout' }"
          @click.stop="handleMenuItemClick(item)"
        >
          <i :class="item.icon"></i>
          <span>{{ item.label }}</span>
        </div>
      </div>

      <!-- Categories container with horizontal scroll -->
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

      <!-- Chat list container -->
      <div class="chat-list-scroll custom-scroll" ref="chatListScroll">
        <div class="chat-list">
          <button
            @click="selectChat('general')"
            :class="{ active: activeChat === 'general' }"
          >
            <i class="bi bi-people"></i> Общий чат
          </button>
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

      <!-- Bottom area for mobile screens -->
      <div class="bottom-area">
        <button class="bi bi-gear"></button>
        <button
          class="theme-toggle"
          @click="emitToggleTheme"
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
        <button class="bi bi-person-circle"></button>
      </div>
    </div>

    <!-- Minimized chat list -->
    <div v-else class="chat-list-minimized">
      <div
        class="minimized-header custom-scroll"
        ref="minimizedHeader"
        @click="toggleMiddlePanel"
      >
        <span class="close-mid-btn"
          ><svg
            class="arrow-icon"
            width="32"
            height="32"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            xmlns="http://www.w3.org/2000/svg"
          >
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </span>
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
</template>

<script>
export default {
  name: "ChatSidebar",
  props: {
    activeChat: {
      type: String,
      default: "",
    },
    isDarkTheme: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      // Panel width from localStorage or default value
      panelWidth: localStorage.getItem("panelWidth")
        ? Number(localStorage.getItem("panelWidth"))
        : 280,
      lastExpandedWidth: 280,
      isResizing: false,
      isSearchActive: false,
      searchQuery: "",
      isMenuOpen: false,
      menuItems: [
        { icon: "bi bi-person-plus", label: "Профиль", key: "profile" },
        { icon: "bi bi-palette", label: "Тема", key: "theme" },
        { icon: "bi bi-gear", label: "Настройки", key: "settings" },
        { icon: "bi bi-box-arrow-left", label: "Выход", key: "logout" },
      ],
      activeCategory: 0,
      categories: ["Все", "Личные", "Группы", "Архив", "Избранное"],
      // Width thresholds for panel resizing
      SNAP_THRESHOLD: 260,
      MIN_THRESHOLD: 120,
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
  computed: {
    // Determine if the middle panel is open based on its width
    isMiddleOpen() {
      return this.panelWidth > this.MIN_THRESHOLD;
    },
  },
  mounted() {
    // Add scroll event listeners for chat list and minimized buttons
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
  },
  beforeUnmount() {
    // Remove global event listeners for resizing and scroll events
    window.removeEventListener("mousemove", this.onResize);
    window.removeEventListener("mouseup", this.stopResize);
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
  },
  methods: {
    // Initialize panel resizing on mousedown event
    initResize(e) {
      e.preventDefault();
      this.isResizing = true;
      document.body.style.userSelect = "none";
      window.addEventListener("mousemove", this.onResize);
      window.addEventListener("mouseup", this.stopResize);
    },

    // Handle resizing of the panel based on mouse movement
    onResize(e) {
      if (!this.isResizing) return;
      const panelRect = this.$refs.middlePanel.getBoundingClientRect();
      const newWidth = e.clientX - panelRect.left;
      if (newWidth < this.MIN_THRESHOLD) {
        this.panelWidth = 80;
      } else {
        this.panelWidth =
          newWidth < this.SNAP_THRESHOLD ? this.SNAP_THRESHOLD : newWidth;
        this.lastExpandedWidth = this.panelWidth;
      }
    },

    // Stop resizing the panel and save the new width to localStorage
    stopResize() {
      if (this.isResizing) {
        this.isResizing = false;
        document.body.style.userSelect = "auto";
        window.removeEventListener("mousemove", this.onResize);
        window.removeEventListener("mouseup", this.stopResize);
        localStorage.setItem("panelWidth", this.panelWidth);
        this.$nextTick(() => {
          // Re-add scroll event listeners if needed
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
      }
    },

    // Handle click on a menu item (e.g., theme toggle)
    handleMenuItemClick(item) {
      if (item.key === "theme") {
        this.emitToggleTheme();
      }
      // Optionally, close the menu after selection
      // this.isMenuOpen = false;
    },

    // Handle horizontal scroll for the categories container
    handleCategoriesScroll(e) {
      const container = this.$refs.categoriesContainer;
      if (container) {
        container.scrollLeft += e.deltaY;
      }
    },

    // Set the active category based on user selection
    selectCategory(index) {
      this.activeCategory = index;
    },

    // Toggle the middle panel open/close state and update width accordingly
    toggleMiddlePanel() {
      if (this.isMiddleOpen) {
        this.lastExpandedWidth = this.panelWidth;
        this.panelWidth = 80;
      } else {
        this.panelWidth = this.lastExpandedWidth || 280;
      }
      localStorage.setItem("panelWidth", this.panelWidth);
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

    // Toggle search input visibility and focus on it when activated
    toggleSearch() {
      this.isSearchActive = !this.isSearchActive;
      if (this.isSearchActive) {
        this.$nextTick(() => {
          this.$refs.searchInput.focus();
        });
      }
    },

    // Toggle dropdown menu visibility
    toggleMenu() {
      this.isMenuOpen = !this.isMenuOpen;
    },

    // Emit event to select a chat
    selectChat(chatType) {
      this.$emit("chat-selected", chatType);
    },

    // Handle scroll event in chat list to add a shadow effect on categories container
    handleChatListScroll() {
      if (!this.$refs.chatListScroll || !this.$refs.categoriesContainer) return;
      const scrollTop = this.$refs.chatListScroll.scrollTop;
      this.$refs.categoriesContainer.classList.toggle(
        "scrolled",
        scrollTop > 0
      );
    },

    // Handle scroll event in minimized buttons to add a shadow effect on header
    handleMinimizedScroll() {
      if (!this.$refs.minimizedButtons) return;
      const scrollTop = this.$refs.minimizedButtons.scrollTop;
      const header = this.$refs.minimizedHeader;
      header?.classList.toggle("scrolled", scrollTop > 0);
      if (header) {
        header.style.boxShadow =
          scrollTop > 0 ? "0 4px 8px -2px rgba(0, 0, 0, 0.25)" : "none";
      }
    },

    // Emit event to toggle the theme
    emitToggleTheme() {
      this.$emit("toggle-theme");
    },
  },
};
</script>

<style scoped>
@import "@/styles/midpanel.css";
</style>
