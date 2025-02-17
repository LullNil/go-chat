<template>
  <!-- Middle panel -->
  <div :class="['middle-panel', { closed: !isMiddleOpen }]">
    <!-- Show middle panel button (appears only on small screens) -->
    <div class="close-mid-btn">
      <span @click="$emit('toggle-middle-panel')" class="bi bi-list"></span>
    </div>

    <!-- If panel open show chat switcher -->
    <div class="chat-switcher" v-if="isMiddleOpen">
      <button
        @click="selectChat('general')"
        :class="{ active: activeChat === 'general' }"
      >
        Общий чат
      </button>

      <button
        @click="selectChat('echo')"
        :class="{ active: activeChat === 'echo' }"
      >
        Эхо-чат
      </button>

      <button v-if="activeChat" @click="closeChat" class="close-btn">
        × Закрыть чат
      </button>

      <!-- Buttons inside the sidebar -->

      <button
        class="bi bi-chat-dots"
        @click="$emit('toggle-middle-panel')"
      ></button>
      <button class="bi bi-gear"></button>

      <!-- Theme toggle button -->
      <button @click="$emit('toggle-theme')">
        <h4 :class="isDarkTheme ? 'bi bi-sun' : 'bi bi-moon'"></h4>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    isMiddleOpen: Boolean,
    isDarkTheme: Boolean,
    activeChat: {
      type: String,
      default: "",
    },
  },
  methods: {
    selectChat(chatType) {
      this.$emit("select-chat", chatType);
    },

    closeChat() {
      this.$emit("close-chat");
    },
  },
};
</script>

<style scoped>
@import "@/styles/midpanel.css";

.chat-switcher {
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

button {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
}

button.active {
  background: #007bff;
  color: white;
}

.close-btn {
  margin-top: 20px;
  background: #dc3545;
  color: white;
}
</style>
