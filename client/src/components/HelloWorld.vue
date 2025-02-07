<template>
  <div :class="['chat-container', { 'dark-theme': isDarkTheme }]">
    <div class="chat-box">
      <div class="chat-title">
        <h2>Chat</h2>
        <button @click="toggleTheme" class="btn-theme">
          <span
            :class="
              isDarkTheme ? 'bi bi-brightness-high-fill' : 'bi bi-moon-fill'
            "
          ></span>
        </button>
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
          v-model="newMessage"
          @keydown.enter.exact.prevent="sendMessage"
          @keydown.enter.ctrl.exact="addNewLine"
          placeholder="Введите сообщение..."
          rows="1"
          ref="textarea"
          @input="adjustTextareaHeight"
        ></textarea>
        <button class="btn-circle" @click="sendMessage">➤</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      ws: null,
      messages: [],
      newMessage: "",
      isDarkTheme: false,
    };
  },
  mounted() {
    this.ws = new WebSocket("ws://localhost:8081/ws");

    this.ws.onmessage = (event) => {
      setTimeout(() => {
        try {
          const msg = JSON.parse(event.data);
          this.messages.push(msg);
          this.$nextTick(this.scrollToBottom);
        } catch (error) {
          console.error("Ошибка парсинга JSON:", error);
        }
      }, 1000);
    };

    this.ws.onopen = () => console.log("WebSocket подключен");
    this.ws.onerror = (error) => console.error("WebSocket ошибка:", error);
    this.ws.onclose = () => console.log("WebSocket закрыт");
  },
  methods: {
    sendMessage() {
      if (!this.newMessage.trim()) return;
      const msg = { user: "Вы", text: this.newMessage };
      this.messages.push(msg);
      this.ws.send(JSON.stringify(msg));
      this.newMessage = "";
      this.$nextTick(() => {
        this.scrollToBottom();
        this.resetTextareaHeight();
      });
    },
    addNewLine(event) {
      event.preventDefault();
      this.newMessage += "\n";
      this.$nextTick(this.adjustTextareaHeight);
    },
    scrollToBottom() {
      const chatContainer = this.$refs.chatContainer;
      if (chatContainer) {
        chatContainer.scrollTop = chatContainer.scrollHeight;
      }
    },
    toggleTheme() {
      this.isDarkTheme = !this.isDarkTheme;
    },
    adjustTextareaHeight() {
      const textarea = this.$refs.textarea;
      textarea.style.height = "auto";
      textarea.style.height = `${Math.min(textarea.scrollHeight, 150)}px`;
    },
    resetTextareaHeight() {
      this.$refs.textarea.style.height = "auto";
    },
  },
};
</script>

<style scoped>
@import "@/styles/chat.css";
</style>
