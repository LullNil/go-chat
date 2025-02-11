import { ref } from "vue";

export default function useWebSocket() {
  const ws = ref(null);

  const connect = (url, messageHandler) => {
    disconnect();
    ws.value = new WebSocket(url);

    ws.value.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        messageHandler(msg);
      } catch (error) {
        console.error("Error parsing message:", error);
      }
    };

    ws.value.onerror = (error) => console.error("WebSocket error:", error);
  };

  const disconnect = () => {
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
  };

  const send = (message) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(message));
    }
  };

  return { connect, disconnect, send };
}
