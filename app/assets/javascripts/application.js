//= require jquery
//= require jquery_ujs
//= require activestorage
//= require rails-ujs
//= require bootstrap-sprockets
//= require_tree .

  // ページ内のすべてのクラスが btn-outline-success である要素を選択
  const buttons = document.querySelectorAll(".btn-outline-success");
  // 選択された要素それぞれに対して、指定された関数を実行
  buttons.forEach(function (button) {
  // ボタン要素から data-url 属性を取得し、その値が現在のページのURL (window.location.pathname) と一致するかを確認
  if (button.getAttribute("data-url") === window.location.pathname) {
    button.classList.add("btn-success"); // `btn-success`クラスを追加
    button.classList.remove("btn-outline-success"); // `btn-outline-success`クラスを取り除く
    }
  });
