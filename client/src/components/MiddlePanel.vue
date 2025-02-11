<template>
  <div :class="['middle-panel', { closed: !isMiddleOpen }]">
    <div class="close-mid-btn">
      <span @click="emit('toggle-middle-panel')" class="bi bi-list"></span>
    </div>

    <div v-if="isMiddleOpen" class="chat-switcher">
      <button
        v-for="chat in chats"
        :key="chat.id"
        @click="selectChat(chat.id)"
        :class="{ active: activeChat === chat.id }"
      >
        {{ chat.label }}
      </button>
      <button v-if="activeChat" @click="closeChat" class="close-btn">
        × Закрыть чат
      </button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from "vue";

// eslint-disable-next-line no-unused-vars
const props = defineProps({
  isMiddleOpen: Boolean,
  activeChat: String,
});

const emit = defineEmits(["toggle-middle-panel", "select-chat", "close-chat"]);

const chats = [
  { id: "general", label: "Общий чат" },
  { id: "echo", label: "Эхо-чат" },
];

const selectChat = (chatType) => emit("select-chat", chatType);
const closeChat = () => emit("close-chat");
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
