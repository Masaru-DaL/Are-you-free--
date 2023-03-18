window.onload = function() {
  drawCanvas();
};

function drawCanvas() {
  // 描画コンテキストの取得
  const canvas = document.getElementById('free-time');
  canvas.style.position = "absolute";
  canvas.style.left = "268px";
  canvas.style.top = "45px"

  if (canvas.getContext) {
    const context = canvas.getContext('2d');
    const lineHeight = 60; // スケジュールの行間隔
    const radius = 15;
    const lineColor = "rgba(195, 128, 41, 0.9)"
    const fillColor = "rgba(212,225,245,0.7)";

    reset();

    function reset() {
      context.clearRect(0,0,canvas.width,canvas.height);

      // 自身のfree-timeを描画する
      for (let i = 0; i < free_times.length; i++) {
        const x = 65 + ((free_times[i].StartHour - 6) * 65) + free_times[i].StartMinute;
        const y = lineHeight;
        const width = ((free_times[i].EndHour - free_times[i].StartHour) * 65) + free_times[i].EndMinute - free_times[i].StartMinute;
        drawFreeTime(x, y, width, lineHeight, radius, lineColor, fillColor);
      }

      // 共有している人のfree-timeを描画する
      for (let i = 0; i < shared_date_free_times.length; i++) {
        const date_free_time = shared_date_free_times[i];
        const free_times = date_free_time.FreeTimes;
        console.log(free_times)

        for (let j = 0; j < free_times.length; j++) {
          const free_time = free_times[j];
          const x = 65 + ((free_time.StartHour - 6) * 65) + free_time.StartMinute;
          const y = lineHeight*(2+i)+((i+1)*20);
          const width = ((free_time.EndHour - free_time.StartHour) * 65) + free_time.EndMinute - free_time.StartMinute;
          drawFreeTime(x, y, width, lineHeight, radius, lineColor, fillColor);
        }
      }
    }

    // free-timeを描画する関数
    function drawFreeTime(x, y, width, height, radius, lineColor, fillColor) {
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
    }
  }
}
