<template>
  <!-- Chat window -->
  <div class="chat-box">
    <!-- If chat selected -->
    <div
      v-if="activeChat"
      class="chat-messages custom-scroll"
      ref="chatMessages"
    >
      <!-- Chat header -->
      <div class="chat-title">
        <button @click="emitToggleSidebar" class="menu-toggle">
          <span class="bi bi-list"></span>
        </button>
        <h2>{{ chatTitle }}</h2>
      </div>

      <div ref="chatContainer" class="chat-messages chat-scroll">
        <!-- Container with dialog start date -->
        <div v-if="messages.length > 0" class="dialog-date">
          {{ dialogStartDate }}
        </div>

        <!-- Chat messages area -->
        <div class="messages-inner">
          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="chat-message"
            :class="msg.user === 'Вы' ? 'sent' : 'received'"
          >
            <p class="chat-text">{{ msg.text }}</p>
          </div>
        </div>
      </div>

      <!-- Message input area -->
      <div class="chat-input">
        <textarea
          class="input-text custom-scroll"
          v-model="newMessage"
          @keydown.enter.exact.prevent="handleSendMessage"
          @keydown.enter.ctrl.exact="addNewLine"
          placeholder="Введите сообщение..."
          rows="1"
          ref="textarea"
          @input="adjustTextareaHeight"
        ></textarea>
        <!-- Send message icon -->
        <button class="btn-circle" @click="handleSendMessage">
          <svg class="send-icon" viewBox="0 0 24 24" fill="none">
            <path
              d="M5.74 15.75L3.5 20.5L20.5 12L3.5 3.5L5.74 8.25L12.5 12L5.74 15.75Z"
              fill="currentColor"
              stroke="currentColor"
              stroke-width="1.2"
              stroke-linejoin="round"
            />
          </svg>
        </button>
      </div>
    </div>

    <!-- If no chat selected -->
    <div v-else class="chat-placeholder">
      <p>Выберите чат для начала общения</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "ChatBox",
  props: {
    activeChat: {
      type: String,
      default: "",
    },
  },

  data() {
    return {
      messages: [],
      ws: null,
      newMessage: "",
      dialogStartDate: new Date().toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "long",
        year: "numeric",
      }),
    };
  },

  computed: {
    chatTitle() {
      switch (this.activeChat) {
        case "general":
          return "Общий чат";
        case "echo":
          return "Эхо-чат";
        default:
          return "Чат";
      }
    },
  },

  watch: {
    activeChat(newChat, oldChat) {
      // Disconnect and clear messages on chat change
      if (oldChat) {
        this.disconnectWebSocket();
        this.messages = [];
      }
      if (newChat) {
        this.connectWebSocket();
      }
    },
  },

  beforeUnmount() {
    this.disconnectWebSocket();
  },

  methods: {
    // Emit toggle-sidebar event
    emitToggleSidebar() {
      this.$emit("toggle-sidebar");
    },

    // Connect to WebSocket
    connectWebSocket() {
      let url = "";
      if (this.activeChat === "general") {
        url = "ws://localhost:8081/ws/general";
      } else if (this.activeChat === "echo") {
        url = "ws://localhost:8081/ws/echo";
      }
      if (url) {
        this.ws = new WebSocket(url);
        this.ws.onopen = () => {
          console.log(`WebSocket connected: ${this.activeChat}`);
        };
        this.ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data);
            if (this.activeChat === "echo") {
              setTimeout(() => {
                this.messages.push(msg);
                this.$nextTick(() => {
                  this.scrollToBottom();
                });
              }, 1000);
            } else {
              this.messages.push(msg);
              this.$nextTick(() => {
                this.scrollToBottom();
              });
            }
          } catch (error) {
            console.error("Error parsing JSON:", error);
          }
        };
        this.ws.onerror = (error) => console.error("WebSocket error:", error);
        this.ws.onclose = () => console.log("WebSocket closed");
      }
    },

    // Disconnect from WebSocket
    disconnectWebSocket() {
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
    },

    // Send message over WebSocket
    handleSendMessage() {
      if (!this.newMessage.trim()) return;
      const msg = { user: "Вы", text: this.newMessage };
      if (this.activeChat === "echo") {
        this.messages.push(msg);
      }
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(msg));
      } else {
        console.error("WebSocket is not open.");
      }
      this.newMessage = "";
      this.$nextTick(() => {
        this.scrollToBottom();
        this.adjustTextareaHeight();
      });
    },

    // Add a new line to the message input
    addNewLine(event) {
      event.preventDefault();
      this.newMessage += "\n";
      this.$nextTick(() => {
        this.adjustTextareaHeight();
        const textarea = this.$refs.textarea;
        if (textarea) {
          textarea.scrollTop = textarea.scrollHeight;
        }
      });
    },

    // Adjust the height of the textarea
    adjustTextareaHeight() {
      const textarea = this.$refs.textarea;
      if (textarea) {
        textarea.style.height = "auto";
        textarea.style.height = `${Math.min(textarea.scrollHeight, 150)}px`;
      }
    },

    // Scroll to the bottom of the chat container
    scrollToBottom() {
      const chatContainer = this.$refs.chatContainer;
      if (chatContainer) {
        const start = chatContainer.scrollTop;
        const end = chatContainer.scrollHeight - chatContainer.clientHeight;
        const change = end - start;
        const duration = 300;
        let startTime = null;
        const animateScroll = (timestamp) => {
          if (!startTime) startTime = timestamp;
          const elapsed = timestamp - startTime;
          const progress = Math.min(elapsed / duration, 1);
          const easedProgress = 1 - Math.pow(1 - progress, 3);
          chatContainer.scrollTop = start + change * easedProgress;
          if (progress < 1) {
            requestAnimationFrame(animateScroll);
          }
        };
        requestAnimationFrame(animateScroll);
      }
    },
  },
};
</script>

<style scoped>
@import "@/styles/chat.css";
</style>
