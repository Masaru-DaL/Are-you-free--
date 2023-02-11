const $icon_list = document.querySelector("#icon-list");
const $icon_list_display_button = document.querySelector("#icon-list-display-button");

$icon_list_display_button.addEventListener("click", function() {
  $icon_list.classList.toggle("is-hidden")
});
