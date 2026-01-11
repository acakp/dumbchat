function scrollToBottom() {
  const chat = document.getElementById('chat');
  chat.scrollTop = chat.scrollHeight;
}
document.addEventListener('DOMContentLoaded', scrollToBottom);

document.body.addEventListener('htmx:afterSwap', function (event) {
  if (event.detail.target.id === 'chat') {
    scrollToBottom();
  }
});
