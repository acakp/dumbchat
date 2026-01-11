// scroll to bottom on new message and page load
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

// send message w/ enter clear msg box after sending
const textarea = document.querySelector('textarea[name="content"]');
const form = document.querySelector('form[hx-post="/messages"]');

if (textarea && form) {
  textarea.addEventListener('keydown', function(e) {
    if (e.key === 'Enter') {
      if (e.ctrlKey || e.metaKey) {
        return // default enter behavior
      } else if (!e.shiftKey){
        // then send
        e.preventDefault();
        if (this.value.trim()) {
          htmx.trigger(form, 'submit');
        }
      }
    }
  });
  
  form.addEventListener('htmx:afterRequest', function() {
    textarea.value = '';
    textarea.focus();
  });

  textarea.focus();
}
