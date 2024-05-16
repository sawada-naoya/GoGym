//= require jquery
//= require activestorage
//= require rails-ujs
//= require bootstrap-sprockets
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


$(document).ready(function() {
  function adjustPadding() {
    var headerHeight = $('header').outerHeight(); // ヘッダーの高さを取得
    $('main').css('padding-top', headerHeight); // mainのpadding-topを設定
    }

  // ページ読み込み時とウィンドウリサイズ時に調整を行う
  adjustPadding();
  $(window).resize(adjustPadding);
});
