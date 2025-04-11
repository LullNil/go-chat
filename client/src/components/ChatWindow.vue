<template>
  <!-- Chat window container -->
  <div class="chat-box">
    <!-- If a chat is active -->
    <div
      v-if="activeChat"
      class="chat-messages custom-scroll"
      ref="chatMessages"
    >
      <!-- Chat header -->
      <div class="chat-title">
        <button @click="closeChatMobile" class="close-chat">
          <span class="close-chat-btn"
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
              <polyline points="15 18 9 12 15 6" />
            </svg>
          </span>
        </button>
        <h2>{{ chatTitle }}</h2>
      </div>

      <div ref="chatContainer" class="chat-messages chat-scroll">
        <!-- Display dialog start date if messages exist -->
        <div v-if="messages.length > 0" class="dialog-date">
          {{ dialogStartDate }}
        </div>

        <!-- Chat messages -->
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

        <!-- Send button -->
        <transition name="puff-send">
          <button
            class="btn-circle send-button"
            v-if="newMessage.trim()"
            @click="handleSendMessage"
            key="send-button"
          >
            <svg
              width="22"
              height="22"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <line
                x1="12"
                y1="19"
                x2="12"
                y2="5"
                stroke="white"
                stroke-width="2"
                stroke-linecap="round"
              />
              <polyline
                points="5,12 12,5 19,12"
                fill="none"
                stroke="white"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
        </transition>
      </div>
    </div>

    <!-- Placeholder if no chat is active -->
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
      // Initialize dialog start date in Russian locale format
      dialogStartDate: new Date().toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "long",
      }),
    };
  },
  computed: {
    // Compute the chat title based on the active chat type
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
    // Watch for changes in activeChat to manage WebSocket connection and message reset
    activeChat(newChat, oldChat) {
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
    // Disconnect WebSocket when component is destroyed
    this.disconnectWebSocket();
  },
  methods: {
    // Emit event to toggle the sidebar visibility
    closeChatMobile() {
      this.$emit("close-chat");
    },

    // Establish WebSocket connection based on the active chat type
    connectWebSocket() {
      let url = "";
      if (this.activeChat === "general") {
        url = "ws://localhost:8081/ws/general";
      } else if (this.activeChat === "echo") {
        url = "ws://localhost:8081/ws/echo";
      }
      if (url) {
        this.ws = new WebSocket(url);
        // On successful connection, log connection status
        this.ws.onopen = () => {
          console.log(`WebSocket connected: ${this.activeChat}`);
        };
        // On receiving a message, parse JSON and update message list
        this.ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data);
            // Delay echo chat messages by 1 second
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

    // Close the WebSocket connection if it exists
    disconnectWebSocket() {
      if (this.ws) {
        this.ws.close();
        this.ws = null;
      }
    },

    // Send message via WebSocket and update message list
    handleSendMessage() {
      if (!this.newMessage.trim()) return;
      const msg = { user: "Вы", text: this.newMessage };
      // For echo chat, immediately add message to the list
      if (this.activeChat === "echo") {
        this.messages.push(msg);
      }
      // Send message if WebSocket is open
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(msg));
      } else {
        console.error("WebSocket is not open.");
      }
      this.newMessage = "";
      this.$nextTick(() => {
        this.$refs.textarea.focus();
      });
      this.$nextTick(() => {
        this.scrollToBottom();
        this.adjustTextareaHeight();
      });
    },

    // Add a newline in the message input when Ctrl+Enter is pressed
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

    // Adjust the textarea height based on content, with a maximum height limit
    adjustTextareaHeight() {
      const textarea = this.$refs.textarea;
      if (textarea) {
        textarea.style.height = "auto";
        textarea.style.height = `${Math.min(textarea.scrollHeight, 150)}px`;
      }
    },

    // Smoothly scroll the chat container to the bottom
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
