import { ref } from "vue";

/**
 * A composable that encapsulates WebSocket connection management.
 * Provides functions to connect, disconnect, and send messages over a WebSocket.
 */
export default function useWebSocket() {
  // Reactive reference to the WebSocket instance
  const ws = ref(null);

  /**
   * Connect to a WebSocket server and set up message handling.
   *
   * @param {string} url - The URL of the WebSocket server.
   * @param {Function} messageHandler - A callback function to process incoming messages.
   */
  const connect = (url, messageHandler) => {
    // Ensure any existing connection is closed before creating a new one
    disconnect();

    // Create a new WebSocket instance with the given URL
    ws.value = new WebSocket(url);

    // Set up the message event handler
    ws.value.onmessage = (event) => {
      try {
        // Parse the incoming data as JSON
        const msg = JSON.parse(event.data);
        // Call the provided message handler with the parsed message
        messageHandler(msg);
      } catch (error) {
        console.error("Error parsing message:", error);
      }
    };

    // Set up the error event handler for the WebSocket
    ws.value.onerror = (error) => console.error("WebSocket error:", error);
  };

  /**
   * Disconnects the current WebSocket connection, if any.
   */
  const disconnect = () => {
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
  };

  /**
   * Sends a message over the WebSocket connection.
   *
   * @param {any} message - The message to send. It will be serialized as JSON.
   */
  const send = (message) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(message));
    }
  };

  // Return the WebSocket API for use in components
  return { connect, disconnect, send };
}
