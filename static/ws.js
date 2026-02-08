function processWsMessage() {
  try {
  const conn = new WebSocket(window.location.origin.replace("http", "ws") + window.chatURLs.ws);
  conn.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      if (msg.type === "new_message") {
        htmx.ajax(
          "GET",
          `${window.chatURLs.message}/${msg.data.id}`,
          { target: "#chat", swap: "beforeend" }
        );
      }

      if (msg.type === "delete_message") {
        const el = document.querySelector(`[data-id="${msg.data.id}"]`);
        if (el) el.remove();
        console.log("msg deleted, id:", msg.data.id)
      }
    };
  } catch(e) {
    console.error('Error processing websocket message:', e);
  }
}

window.onload = processWsMessage();
