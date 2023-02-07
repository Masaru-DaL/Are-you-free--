document.getElementById("today").innerHTML = getNow();

function getNow() {
  var now = new Date();
  var year = now.getFullYear();
  var month = now.getMonth()+1;
  var day = now.getDate();
  var hour = now.getHours();
  var minute = now.getMinutes();
  var nowDayOfWeek = now.getDay();

  // 曜日
  var dayOfWeek = new Array("日","月","火","水","木","金","土");

  var nowDate = year + "/" + month + "/" + day + "（" + dayOfWeek[nowDayOfWeek] + "）" + hour + "：" + minute

  return nowDate;
}
