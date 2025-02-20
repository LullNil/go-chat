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
      <button @click="$emit('toggle-theme')">
        <span :class="isDarkTheme ? 'bi bi-sun' : 'bi bi-moon'"></span>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showSettings: false,
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
