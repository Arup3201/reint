(function () {
  function generateData() {
    let timestamps = [];
    let actual = [];
    let forecast = [];

    let current = new Date("2024-01-01T00:00:00");
    let end = new Date("2024-01-31T23:30:00");

    let i = 0;

    while (current <= end) {
      timestamps.push(new Date(current));

      let base = 6000 + Math.sin(i / 15) * 1500;

      actual.push([current, base + Math.random() * 300]);
      forecast.push([current, base + Math.random() * 500 - 200]);

      current = new Date(current.getTime() + 30 * 60 * 1000);

      i++;
    }

    return { timestamps, actual, forecast };
  }

  const dataset = generateData();

  const chart = echarts.init(document.getElementById("chart"));

  const option = {
    tooltip: {
      trigger: "axis",
    },
    legend: {
      data: ["Actual", "Forecast"],
    },
    xAxis: {
      type: "time",
      name: "Timestamp (Jan 2024)",
    },
    yAxis: {
      type: "value",
      name: "Wind Power (MW)",
    },
    dataZoom: [
      {
        type: "inside",
      },
      {
        type: "slider",
      },
    ],
    series: [
      {
        name: "Actual",
        type: "line",
        data: dataset.actual,
        smooth: true,
        lineStyle: { color: "blue" },
      },
      {
        name: "Forecast",
        type: "line",
        data: dataset.forecast,
        smooth: true,
        lineStyle: { color: "green" },
      },
    ],
  };

  chart.setOption(option);

  // slider update
  const slider = document.getElementById("horizonSlider");
  const horizonText = document.getElementById("horizonValue");

  slider.oninput = () => {
    horizonText.innerText = slider.value + " Hr";
  };

  // resize responsiveness
  window.addEventListener("resize", () => {
    chart.resize();
  });
})();
