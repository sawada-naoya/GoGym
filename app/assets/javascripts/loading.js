document.addEventListener('DOMContentLoaded', function() {
  var loader = document.querySelector('.loader');
  console.log(loader); // これでloader要素が正しく取得されているか確認

  if (loader) {
    // ページの読み込みが完了したらアニメーションを非表示
    window.addEventListener('load', function() {
      console.log('Page loaded'); // これが表示されるか確認
      loader.style.display = 'none';
    });

    // ページの読み込みが完了してなくても2秒後にアニメーションを非表示にする
    setTimeout(function() {
      console.log('Timeout reached'); // これが表示されるか確認
      loader.style.display = 'none';
    }, 2000);
  }
});
