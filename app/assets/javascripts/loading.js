document.addEventListener('DOMContentLoaded', function() {
  var loader = document.querySelector('.loader');

  // ページの読み込みが完了したらアニメーションを非表示
  window.addEventListener('load', function() {
    loader.style.display = 'none';
  });

  // ページの読み込みが完了してなくても2秒後にアニメーションを非表示にする
  setTimeout(function() {
    loader.style.display = 'none';
  }, 2000);
});
