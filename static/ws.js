function processWsMessage() {
  try {
  const conn = new WebSocket(window.location.origin.replace("http", "ws") + "/chat/ws");
  conn.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      if (msg.type === "new_message") {
        htmx.ajax(
          "GET",
          `/chat/message/${msg.data.id}`,
          { swap: "beforeend" }
        );
      }

      if (msg.type === "delete_message") {
        const el = document.querySelector(`[data-id="${msg.data}"]`);
        if (el) el.remove();
      }
    };
  } catch(e) {
    console.error('Error processing websocket message:', e);
  }

  // code below is deprecated
  // try {
  //   const data = JSON.parse(event.data);
  //   if (data.type === 'new_message') {
  //     // create a temporary div to hold the message html using the same format as the template
  //     const tempDiv = document.createElement('div');
  //     tempDiv.innerHTML = `
  //       <div class="message" data-id="${data.data.id}">
  //         <div class="message-line">
  //           <span class="time">${data.data.formattedTime}</span>
  //           <span class="sender">${data.data.nickname}:</span>
  //           <span class="text">${data.data.content}</span>
  //         </div>
  //       </div>
  //     `;

  //     const chatContainer = document.getElementById('chat');
  //     chatContainer.appendChild(tempDiv.firstElementChild);

  //     chatContainer.scrollTop = chatContainer.scrollHeight;
  //   }
  // } catch (e) {
  //   console.error('Error processing websocket message:', e);
  // }
}

window.onload = processWsMessage();
