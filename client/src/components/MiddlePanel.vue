<template>
  <div :class="['middle-panel', { closed: !isMiddleOpen }]">
    <!-- Кнопка для открытия/закрытия панели -->
    <div class="close-mid-btn">
      <span @click="$emit('toggle-middle-panel')" class="bi bi-list"></span>
    </div>

    <!-- Полная версия списка чатов -->
    <div class="chat-list" v-if="isMiddleOpen">
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
      <!-- Кнопка закрытия чата удалена -->
    </div>

    <!-- Минимизированная версия списка чатов (при закрытой панели) -->
    <div class="chat-list-minimized" v-else>
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

    <!-- Нижняя область для настроек и смены темы -->
    <div class="bottom-area" v-if="isMiddleOpen">
      <button class="bi bi-gear"></button>
      <button class="bi bi-person-circle"></button>
      <button
        class="theme-toggle"
        @click="$emit('toggle-theme')"
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
</template>

<script>
export default {
  data() {
    return {
      // Sun icon
      sunIcon: `M12 3v1
                m0 16v1
                m9-9h-1
                M4 12H3
                m15.364-6.364l-.707.707
                M6.343 17.657l-.707.707
                m12.728 0l-.707-.707
                M6.343 6.343l-.707-.707
                M12 8a4 4 0 110 8 4 4 0 010-8z`,
      // Moon icon
      moonIcon: `M21 12.79A9 9 0 1111.21 3
                 a7 7 0 009.79 9.79z`,
      showSettings: false,
      animateIcon: false,
    };
  },
  props: {
    isMiddleOpen: Boolean,
    isDarkTheme: Boolean,
    activeChat: {
      type: String,
      default: "",
    },
  },

  mounted() {
    // Слушатель для кнопок мыши. Значение e.button === 4 соответствует пятой кнопке (обычно "Forward").
    window.addEventListener("mouseup", this.handleMouseButton);
    // Слушатель для сочетания клавиш. Здесь пример: Ctrl+Shift+X
    window.addEventListener("keydown", this.handleKeyCombination);
  },
  beforeUnmount() {
    window.removeEventListener("mouseup", this.handleMouseButton);
    window.removeEventListener("keydown", this.handleKeyCombination);
  },
  methods: {
    selectChat(chatType) {
      this.$emit("select-chat", chatType);
    },
    closeChat() {
      this.$emit("close-chat");
    },
    handleMouseButton(e) {
      // e.button: 0 - левая, 1 - средняя, 2 - правая, 3 - кнопка "Назад", 4 - кнопка "Вперёд"
      if (e.button === 3) {
        e.preventDefault(); // предотвращаем возможное действие по умолчанию (например, переход вперёд)
        this.closeChat();
      }
    },
    handleKeyCombination(e) {
      // Пример: Ctrl+Shift+X для закрытия чата
      if (e.ctrlKey && e.shiftKey && e.code === "KeyX") {
        e.preventDefault();
        this.closeChat();
      }
    },
  },
};
</script>

<style scoped>
@import "@/styles/midpanel.css";
</style>
