<template>
  <!-- Sidebar Panel for chat navigation -->
  <div
    ref="sidebarPanel"
    class="sidebar-panel"
    :class="{ closed: !isSidebarOpen, resizing: isResizing }"
    :style="{ width: isSidebarOpen ? panelWidth + 'px' : '80px' }"
  >
    <!-- Resize handle for Sidebar panel -->
    <div
      class="resize-handle"
      :class="{ active: isResizing }"
      @mousedown="initResize"
    ></div>

    <!-- Expanded Sidebar Panel Content -->
    <div v-if="isSidebarOpen" class="sidebar-panel-list">
      <!-- Sidebar header with dropdown menu trigger -->
      <div class="sidebar-header">
        <div
          class="plus-icon-wrapper"
          :class="{ active: isMenuOpen }"
          @click="toggleMenu"
        >
          <svg class="plus-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
          </svg>
        </div>
        <div class="sidebar-header-action">
          <h3
            @click="toggleSidebarPanel"
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
                d="M15.8 16.2L20 20.4L18.6 21.8L14.4 17.6C13.2 18.5 11.7 19 10 19C5.6 19 2 15.4 2 11C2 6.6 5.6 3 10 3C14.4 3 18 6.6 18 11C18 12.7 17.5 14.2 16.6 15.4L16.2 15.8H15.8ZM10 17C13.3 17 16 14.3 16 11C16 7.7 13.3 5 10 5C6.7 5 4 7.7 4 11C4 14.3 6.7 17 10 17Z"
              />
            </svg>
          </div>
        </div>
      </div>

      <!-- Dropdown menu -->
      <div class="dropdown-menu-mid" :class="{ active: isMenuOpen }">
        <div v-for="(item, index) in menuItems" :key="index">
          <!-- Divider line -->
          <div v-if="item.hasDivider" class="divider"></div>
          <div
            class="menu-item"
            :class="{ 'logout-item': item.key === 'logout' }"
            @click.stop="handleMenuItemClick(item)"
          >
            <i :class="item.icon"></i>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </div>

      <!-- Categories container with horizontal scroll -->
      <div class="categories-container" ref="categoriesContainer">
        <div
          class="categories-scroll"
          ref="categoriesScroll"
          @wheel.passive="handleCategoriesScroll"
        >
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
        <button class="bi bi-person-circle"></button>
      </div>
    </div>

    <!-- Minimized chat list -->
    <div v-else class="chat-list-minimized">
      <div
        class="minimized-header custom-scroll"
        ref="minimizedHeader"
        @click="toggleSidebarPanel"
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
import { mapGetters, mapActions } from "vuex";

export default {
  name: "ChatSidebar",
  props: {
    activeChat: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      // Panel width from localStorage or default value
      panelWidth: localStorage.getItem("panelWidth")
        ? Number(localStorage.getItem("panelWidth"))
        : 260,
      lastExpandedWidth: 260,
      isResizing: false,
      isSearchActive: false,
      searchQuery: "",
      isMenuOpen: false,
      menuItems: [
        { icon: "bi bi-person-plus", label: "Профиль", key: "profile" },
        {
          icon: "bi bi-palette",
          label: "Тема",
          key: "theme",
          hasDivider: true,
        },
        { icon: "bi bi-gear", label: "Настройки", key: "settings" },
        {
          icon: "bi bi-box-arrow-left",
          label: "Выход",
          key: "logout",
          hasDivider: true,
        },
      ],
      activeCategory: 0,
      categories: ["Все", "Непрочитанные", "Группы", "Избранное", "Архив"],
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

  mounted() {
    // Add scroll event listeners for chat list and minimized buttons
    if (this.$refs.chatListScroll) {
      this.$refs.chatListScroll.addEventListener(
        "scroll",
        this.handleChatListScroll,
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
    const isDark = this.$store.getters["theme/isDarkTheme"];
    if (isDark) {
      document.documentElement.classList.add("dark-theme");
    } else {
      document.documentElement.classList.remove("dark-theme");
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

  computed: {
    // Determine if the sidebar panel is open based on its width
    isSidebarOpen() {
      return this.panelWidth > this.MIN_THRESHOLD;
    },

    ...mapGetters("auth", ["isAuthenticated"]),
    ...mapGetters("theme", ["isDarkTheme"]),
  },

  methods: {
    ...mapActions("auth", ["logout"]),
    ...mapActions("theme", ["toggleTheme"]),

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
      const panelRect = this.$refs.sidebarPanel.getBoundingClientRect();
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

    // Handle click on a dropdown-menu item
    handleMenuItemClick(item) {
      switch (item.key) {
        case "theme":
          this.toggleTheme();
          break;
        case "logout":
          this.$store.dispatch("auth/logout");
          this.$router.push({ name: "AuthPage" });
          break;
        case "profile":
          // console.log("Go to profile");
          break;
        case "settings":
          // console.log("Open settings");
          break;
        default:
          console.log("Dropdown item error:", item.key);
      }
      // this.isMenuOpen = false;
    },

    // Handle horizontal scroll for the categories container
    handleCategoriesScroll(e) {
      const container = this.$refs.categoriesScroll;
      if (container) {
        const delta = e.deltaY || e.deltaX;
        container.scrollLeft += delta * 0.3;
      }
    },

    // Set the active category based on user selection
    selectCategory(index) {
      this.activeCategory = index;
    },

    // Toggle the sidebar panel open/close state and update width accordingly
    toggleSidebarPanel() {
      if (this.isSidebarOpen) {
        this.lastExpandedWidth = this.panelWidth;
        this.panelWidth = 80;
      } else {
        this.panelWidth = this.lastExpandedWidth || 260;
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
  },
};
</script>

<style scoped>
@import "@/styles/sidebar.css";
</style>
