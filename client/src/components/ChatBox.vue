<template>
  <div class="chat-box">
    <div class="chat-title">
      <!-- Button that appears only on small screens -->
      <button @click="$emit('toggle-middle-panel')" class="menu-toggle">
        <span class="bi bi-list"></span>
      </button>
      <h2>Chat</h2>
    </div>

    <!-- Chat content -->
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

    <!-- Input field and send button -->
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

      <!-- Send button -->
      <button class="btn-circle" @click="sendMessage">➤</button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    messages: Array,
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
        this.$refs.chatContainer.scrollTop =
          this.$refs.chatContainer.scrollHeight;
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
