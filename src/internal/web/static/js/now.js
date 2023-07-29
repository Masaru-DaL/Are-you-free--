// function getNowInput () {
//   var today = new Date();
//   today.setDate(today.getDate());
//   var yyyy = today.getFullYear();
//   var mm = ("0"+(today.getMonth()+1)).slice(-2);
//   var dd = ("0"+today.getDate()).slice(-2);
//   document.getElementById("today-input").value=yyyy+'-'+mm+'-'+dd;
// }

function set2fig(num) {
  // 桁数が1桁だったら先頭に0を加えて2桁に調整する
  var ret;
  if( num < 10 ) { ret = "0" + num; }
  else { ret = num; }
  return ret;
}

function getNow() {
  var now = new Date();
  var year = now.getFullYear();
  var month = now.getMonth()+1;
  var day = now.getDate();
  var hour = set2fig(now.getHours());
  var minute = set2fig(now.getMinutes());
  var second = set2fig(now.getSeconds());
  var nowDayOfWeek = now.getDay();

  // 曜日
  var dayOfWeek = new Array("日","月","火","水","木","金","土");

  var nowDate = year + "/" + month + "/" + day + "（" + dayOfWeek[nowDayOfWeek] + "）" + hour + "：" + minute + "：" + second;

  document.getElementById("today").innerHTML = nowDate;
}
document.addEventListener('DOMContentLoaded', function() {
  setInterval(getNow, 1000);

  var todayDiv = document.getElementById('today-wrapper');
  todayDiv.style.margin = '0 auto';
  todayDiv.style.padding= '0.5% 1%';
  todayDiv.style.border= 'none';
  todayDiv.style.borderRight= '1.5px solid var(--WOOD)';
  todayDiv.style.borderBottom= '1.5px solid var(--WOOD)';
  todayDiv.style.boxShadow= '7px 7px 7px rgba(54, 54, 54, 0.3)';
  todayDiv.style.backgroundColor= 'var(--SNOW)';


  // free_time.forEach(function(element) {
  //   console.log("ID: " + element.ID);
  //   console.log("DateFreeTimeID: " + element.DateFreeTimeID);
  //   console.log("StartHour: " + element.StartHour);
  //   console.log("StartMinute: " + element.StartMinute);
  //   console.log("EndHour: " + element.EndHour);
  //   console.log("EndMinute: " + element.EndMinute);
  //   console.log("CreatedAt: " + element.CreatedAt);
  //   console.log("UpdatedAt: " + element.UpdatedAt);
  // });
});
// setInterval(getNow, 1000);
// getNowInput();
