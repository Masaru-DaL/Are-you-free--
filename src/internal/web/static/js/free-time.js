window.onload = function() {
  drawCanvas();
};

// function drawCanvas() {
//   console.log(free_times);

//   // 描画コンテキストの取得
//   const canvas = document.getElementById('free-time');
//   canvas.style.position = "absolute";
//   canvas.style.left = "270px";
//   canvas.style.top = "150px"

//   if (canvas.getContext) {
//     const context = canvas.getContext('2d');

//     // free-scheduleの変数定義
//     const x = 50 + ((free_times[0].StartHour-6) * 60) + free_times[0].StartMinute;
//     const y = 1 * 70;
//     const width = ((free_times[0].EndHour - free_times[0].StartHour) * 60) + free_times[0].endminute - free_times[0].startminute;
//     // const x = 50 + (({{index . "starthour"free_times[0].StartHour}}-6) * 60) + {{index . "startminute"}};
//     // const y = {{index . "id"}} * 70;
//     // const width = (({{index . "endhour"}} - {{index . "starthour"}}) * 60) + {{index . "endminute"}} - {{index . "startminute"}};
//     const height = 70;
//     const radius = 15;
//     const lineColor = "rgba(215,155,0.1)"
//     const fillColor = "rgba(212,225,245,0.4)";
//     const color = "rgba(212,225,245,1)";

//     reset();

//     function reset() {
//       context.clearRect(0,0,canvas.width,canvas.height);
//       drawSchedule(x, y, width, height, radius, lineColor, fillColor);
//     }

//     // スケジュールを描画する関数（角を丸める）
//     function drawSchedule(x, y, width, height, radius, lineColor, fillColor) {
//       context.beginPath();
//       context.lineWidth = 1;
//       context.strokeStyle = lineColor;
//       context.fillStyle = fillColor;
//       context.moveTo(x,y + radius);
//       context.arc(x+radius, y+height-radius, radius, Math.PI, Math.PI*0.5, true);
//       context.arc(x+width-radius, y+height-radius, radius, Math.PI*0.5, 0,1);
//       context.arc(x+width-radius, y+radius, radius, 0, Math.PI*1.5, 1);
//       context.arc(x+radius, y+radius, radius, Math.PI*1.5, Math.PI, 1);
//       context.closePath();
//       context.stroke();
//       context.fill();
//       console.log(free_times);
//     }
//   }
// }

function drawCanvas() {
  // 描画コンテキストの取得
  const canvas = document.getElementById('free-time');
  canvas.style.position = "absolute";
  canvas.style.left = "260px";
  canvas.style.top = "30px"

  if (canvas.getContext) {
    const context = canvas.getContext('2d');
    const lineHeight = 60; // スケジュールの行間隔
    const radius = 15;
    const lineColor = "rgba(195, 128, 41, 0.9)"
    // const fillColor = "rgba(0,100,0,0.3)";
    // const fillColor = "rgba(0,100,0,0.3)";
    const fillColor = "rgba(212,225,245,0.7)";

    reset();

    function reset() {
      context.clearRect(0,0,canvas.width,canvas.height);
      for (let i = 0; i < free_times.length; i++) {
        console.log(free_times[i].StartHour)
        console.log(free_times[i].StartMinute)
        console.log(free_times[i].EndHour)
        console.log(free_times[i].EndMinute)

        const x = 65 + ((free_times[i].StartHour - 6) * 65) + free_times[i].StartMinute;
        const y = lineHeight;
        const width = ((free_times[i].EndHour - free_times[i].StartHour) * 65) + free_times[i].EndMinute - free_times[i].StartMinute;
        drawSchedule(x, y, width, lineHeight, radius, lineColor, fillColor);
      }
    }

    matchField();





    // スケジュールを描画する関数（角を丸める）
    function drawSchedule(x, y, width, height, radius, lineColor, fillColor) {
      context.beginPath();
      context.lineWidth = 1;
      context.strokeStyle = lineColor;
      context.fillStyle = fillColor;
      context.moveTo(x,y + radius);
      context.arc(x+radius, y+height-radius, radius, Math.PI, Math.PI*0.5, true);
      context.arc(x+width-radius, y+height-radius, radius, Math.PI*0.5, 0,1);
      context.arc(x+width-radius, y+radius, radius, 0, Math.PI*1.5, 1);
      context.arc(x+radius, y+radius, radius, Math.PI*1.5, Math.PI, 1);
      context.closePath();
      context.stroke();
      context.fill();
      console.log(free_times);

    }
  }
}
