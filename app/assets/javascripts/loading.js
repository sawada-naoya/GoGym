document.addEventListener('turbo:load', function() {
  var loaderWrap = document.querySelector('.loader-wrap');
  console.log("ローダー要素を取得:", loaderWrap); // ローダー要素が正しく取得されているか確認

  if (loaderWrap) {
    loaderWrap.classList.remove('loading'); // 初期表示を非表示に設定
  }
});

document.addEventListener('turbo:request-start', function() {
  var loaderWrap = document.querySelector('.loader-wrap');
  if (loaderWrap) {
    console.log("リクエスト開始"); // これが表示されるか確認
    loaderWrap.classList.add('loading'); // リクエスト開始時にローダーを表示
  }
});

document.addEventListener('turbo:request-end', function() {
  var loaderWrap = document.querySelector('.loader-wrap');
  if (loaderWrap) {
    console.log("リクエスト終了"); // これが表示されるか確認
    loaderWrap.classList.remove('loading'); // リクエスト終了時にローダーを非表示
  }
});

document.addEventListener('DOMContentLoaded', function() {
  var loaderWrap = document.querySelector('.loader-wrap');
  console.log("ローダー要素を取得:", loaderWrap); // ローダー要素が正しく取得されているか確認

  if (loaderWrap) {
    // ページの読み込みが完了したらアニメーションを非表示
    window.addEventListener('load', function() {
      console.log("ページ読み込み完了"); // これが表示されるか確認
      loaderWrap.classList.remove('loading');
    });

    // ページの読み込みが完了してなくても2秒後にアニメーションを非表示にする
    setTimeout(function() {
      console.log("タイムアウト到達"); // これが表示されるか確認
      loaderWrap.classList.remove('loading');
    }, 2000);
  }
});
