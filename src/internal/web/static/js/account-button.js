/* icon変更ボタン */
const icon_list = document.querySelector('#icon-list');
const icon_list_display_button = document.querySelector("#icon-list-display-button");

/* username変更ボタン */
const username_change_form = document.querySelector('#username-change-form');
const username_change_form_display_button = document.querySelector('#username-change-form-display-button');

/* password変更ボタン */
const password_change_form = document.querySelector('#password-change-form');
const password_change_form_display_button = document.querySelector('#password-change-form-display-button');

icon_list_display_button.addEventListener("click", () => {
  icon_list.classList.toggle("icon-list-is-hidden");

  if (username_change_form.classList.contains('username-change-form-is-hidden')) {
    username_change_form.classList.remove("username-change-form-is-hidden");
  };

  if (password_change_form.classList.contains('password-change-form-is-hidden')) {
    password_change_form.classList.remove("password-change-form-is-hidden");
  };
});

username_change_form_display_button.addEventListener("click", () => {
  username_change_form.classList.toggle("username-change-form-is-hidden");

  if (icon_list.classList.contains('icon-list-is-hidden')) {
    icon_list.classList.remove("icon-list-is-hidden");
  };

  if (password_change_form.classList.contains('password-change-form-is-hidden')) {
    password_change_form.classList.remove("password-change-form-is-hidden");
  };
});

password_change_form_display_button.addEventListener("click", () => {
  password_change_form.classList.toggle("password-change-form-is-hidden");

  if (icon_list.classList.contains('icon-list-is-hidden')) {
    icon_list.classList.classList.remove("icon-list-is-hidden");
  };

  if (username_change_form.classList.contains('username-change-form-is-hidden')) {
    username_change_form.classList.remove("username-change-form-is-hidden");
  };
});

// // ボタンのDOM要素を取得
// var three_type_btn = document.getElementsByClassName('title-area-button-in-account');

// // ボタンの個数分ループ
// // 変数「i」に現在のループ回数が代入される
// for (var i = three_type_btn.length -1; i >= 0; i--) {
//   btnAction(three_type_btn[i],i);
// };

// function btnAction(three_type_btnDOM,three_type_btnId){
//   // 各ボタンをイベントリスナーに登録
//   three_type_btnDOM.addEventListener("click", function(){
//     // activeクラスの追加と削除
//     // thisは、クリックされたオブジェクト
//     this.classList.toggle('icon-list-wrapper');

//     // クリックされていないボタンにactiveがついていたら外す
//     for (var i = three_type_btn.length - 1; i >= 0; i--) {
//       if(three_type_btnId !== i){
//         if(three_type_btn[i].classList.contains('icon-list-wrapper')){
//           three_type_btn[i].classList.remove('icon-list-wrapper');
//         }
//       }
//     }
//   });
// };

// ボタンのDOM要素を取得
var btn = document.getElementsByClassName('icon');
var icon_value = document.getElementsByClassName('icon').value;

// ボタンの個数分ループ
// 変数「i」に現在のループ回数が代入される
for (var i = btn.length - 1; i >= 0; i--) {
  btnAction(btn[i],i);
};

function btnAction(btnDOM,btnId){
  // 各ボタンをイベントリスナーに登録
  btnDOM.addEventListener("click", function(){
    // activeクラスの追加と削除
    // thisは、クリックされたオブジェクト
    this.classList.toggle('scale');
    icon_value = this.value

    // クリックされていないボタンにactiveがついていたら外す
    for (var i = btn.length - 1; i >= 0; i--) {
      if(btnId !== i){
        if(btn[i].classList.contains('scale')){
          btn[i].classList.remove('scale');
        }
      }
    }
    document.getElementById('icon-submit-button').onclick = () => {
          console.log(icon_value);
    };
  });
};



// document.querySelectorAll('.icon').forEach(function (button) {
//   button.addEventListener('click', {value: `${icon.value}`, handleEvent: onClickButton});
// });

// function onClickButton() {
//   icon_value = this.value;
//   $icon.classList.toggle('scale');
//   console.log(icon_value);
//   document.getElementById('submit-button').onclick = () => {
//     console.log(icon_value);
//   };
// };

// $icon.addEventListener('click', () => {
// });

// document.getElementById('submit-button').onclick = () => {
//   console.log(Element.classList.value)
// };
