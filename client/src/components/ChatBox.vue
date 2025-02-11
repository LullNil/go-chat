<template>
  <div class="chat-box">
    <!-- Если выбран чат, отображаем его содержимое -->
    <div v-if="activeChat" class="chat-messages custom-scroll">
      <div class="chat-title">
        <!-- Кнопка для маленьких экранов (можно оставить для переключения панели) -->
        <button @click="$emit('toggle-middle-panel')" class="menu-toggle">
          <span class="bi bi-list"></span>
        </button>
        <h2>Chat</h2>
      </div>

      <!-- Область для сообщений -->
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

      <!-- Поле ввода сообщения и кнопка отправки -->
      <div class="chat-input">
        <textarea
          class="input-text custom-scroll"
          v-model="newMessage"
          @keydown.enter.exact.prevent="sendMessage()"
          @keydown.enter.ctrl.exact="addNewLine"
          placeholder="Введите сообщение..."
          rows="1"
          ref="textarea"
          @input="adjustTextareaHeight"
        ></textarea>
        <button class="btn-circle" @click="sendMessage()">➤</button>
      </div>
    </div>
    <!-- Если чат не выбран, показываем placeholder -->
    <div v-else class="chat-placeholder">
      <p>Выберите чат для начала общения</p>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    activeChat: {
      type: String,
      default: "",
    },

    messages: {
      type: Array,
      default: () => [],
    },
  },

  data() {
    return {
      newMessage: "",
    };
  },

  watch: {
    messages: {
      handler() {
        this.$nextTick(() => {
          const chatContainer = this.$refs.chatContainer;
          if (chatContainer) {
            chatContainer.scrollTop = chatContainer.scrollHeight;
          }
        });
      },
      deep: true,
    },
  },

  methods: {
    sendMessage() {
      if (!this.newMessage.trim()) return;
      const msg = { user: "Вы", text: this.newMessage };
      this.$emit("send-message", msg);
      this.newMessage = "";
      this.$nextTick(() => {
        const chatContainer = this.$refs.chatContainer;
        if (chatContainer) {
          chatContainer.scrollTop = chatContainer.scrollHeight;
        }
        this.adjustTextareaHeight();
      });
    },

    addNewLine(event) {
      event.preventDefault();
      this.newMessage += "\n";
    },

    adjustTextareaHeight() {
      const textarea = this.$refs.textarea;
      textarea.style.height = "auto";
      textarea.style.height = `${Math.min(textarea.scrollHeight, 150)}px`;
    },
  },
};
</script>

<style scoped>
@import "@/styles/chat.css";
</style>
