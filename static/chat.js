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

textarea = document.querySelector('textarea[name="content"]');
form = document.querySelector('div.input-area form');

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

// track last msg id
let lastMessageId = 0;

document.body.addEventListener("htmx:afterSwap", function (e) {
  if (e.target.id === "chat") {
    const msgs = e.target.querySelectorAll(".message");
    if (msgs.length > 0) {
      lastMessageId = msgs[msgs.length - 1].dataset.id;
    }
  }
});

// document.body.addEventListener("htmx:configRequest", function (e) {
//   if (e.detail.path === "/poll") {
//     e.detail.parameters.after_id = lastMessageId;
//   }
// });
