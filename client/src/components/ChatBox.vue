<template>
  <div class="chat-box">
    <div v-if="activeChat" class="chat-messages custom-scroll">
      <div class="chat-title">
        <button @click="$emit('toggle-middle-panel')" class="menu-toggle">
          <span class="bi bi-list"></span>
        </button>
        <!-- Заголовок чата теперь зависит от выбранного чата -->
        <h2>{{ chatTitle }}</h2>
      </div>

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

      <div class="chat-input">
        <textarea
          class="input-text custom-scroll"
          ref="textarea"
          v-model="newMessage"
          @keydown.enter.exact.prevent="sendMessage"
          @keydown.enter.ctrl.exact="addNewLine"
          placeholder="Введите сообщение..."
          rows="1"
          @input="adjustHeight"
        ></textarea>
        <button class="btn-circle" @click="sendMessage">➤</button>
      </div>
    </div>

    <div v-else class="chat-placeholder">
      <p>Выберите чат для начала общения</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, computed, defineProps, defineEmits } from "vue";

const props = defineProps({
  activeChat: String,
  messages: Array,
});

const emit = defineEmits(["send-message", "toggle-middle-panel"]);

const newMessage = ref("");
const textarea = ref(null);
const chatContainer = ref(null);

watch(
  () => props.messages,
  () => {
    nextTick(() => {
      if (chatContainer.value) {
        chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
        adjustHeight();
      }
    });
  },
  { deep: true }
);

const sendMessage = () => {
  if (!newMessage.value.trim()) return;
  emit("send-message", { user: "Вы", text: newMessage.value.trim() });
  newMessage.value = "";
};

const addNewLine = (e) => {
  e.preventDefault();
  newMessage.value += "\n";
};

const adjustHeight = () => {
  if (textarea.value) {
    textarea.value.style.height = "auto";
    textarea.value.style.height = `${Math.min(
      textarea.value.scrollHeight,
      150
    )}px`;
  }
};

// Вычисляемое свойство для определения заголовка чата по значению activeChat
const chatTitle = computed(() => {
  if (props.activeChat === "general") {
    return "Общий чат";
  } else if (props.activeChat === "echo") {
    return "Эхо-чат";
  }
  return "Chat";
});
</script>

<style scoped>
@import "@/styles/chat.css";
</style>
