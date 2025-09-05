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

// document.addEventListener('DOMContentLoaded', function() {
//   function adjustMainPadding() {
//     const navbarHeight = document.querySelector('.navbar').offsetHeight;
//     document.querySelector('main').style.paddingTop = `${navbarHeight}px`;
//   }

//   // Adjust padding on page load
//   adjustMainPadding();

//   // Adjust padding on window resize
//   window.addEventListener('resize', adjustMainPadding);
// });
