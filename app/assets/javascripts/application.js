//= require jquery
//= require activestorage
//= require rails-ujs
//= require bootstrap-sprockets
//= require raty
//= require_tree .

// ページ内のすべてのクラスが btn-outline-success である要素を選択
document.addEventListener("DOMContentLoaded", function() {
  const buttons = document.querySelectorAll(".btn-outline-success");
  buttons.forEach(function(button) {
    if (button.getAttribute("data-url") === window.location.pathname) {
      button.classList.add("btn-success");
      button.classList.remove("btn-outline-success");
    }
  });
});
