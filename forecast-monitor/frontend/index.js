(function () {
  const chart = echarts.init(document.getElementById("chart"));

  const startInput = document.getElementById("startTime");
  const endInput = document.getElementById("endTime");
  const slider = document.getElementById("horizonSlider");
  const horizonText = document.getElementById("horizonValue");

  // ── Chart base option (set once) ───────────────────────────────────────────
  const option = {
    tooltip: { trigger: "axis" },
    legend: { data: ["Actual", "Forecast"] },
    xAxis: {
      type: "time",
      name: "Timestamp (Jan 2024)",
    },
    yAxis: {
      type: "value",
      name: "Wind Power (MW)",
    },
    dataZoom: [{ type: "inside" }, { type: "slider" }],
    series: [
      {
        name: "Actual",
        type: "line",
        data: [],
        smooth: true,
        lineStyle: { color: "blue" },
      },
      {
        name: "Forecast",
        type: "line",
        data: [],
        smooth: true,
        lineStyle: { color: "green" },
      },
    ],
  };

  chart.setOption(option);

  // ── Helpers ────────────────────────────────────────────────────────────────

  /** Convert datetime-local value (YYYY-MM-DDTHH:mm) → ISO-8601 UTC string */
  function toISOString(localValue) {
    // datetime-local gives "YYYY-MM-DDTHH:mm"; append ":00Z" for UTC
    return localValue + ":00Z";
  }

  /** Build the API URL from current control values */
  function buildURL() {
    const start = toISOString(startInput.value);
    const end = toISOString(endInput.value);
    const horizon = slider.value;
    return `/wind-data?start=${encodeURIComponent(start)}&end=${encodeURIComponent(end)}&horizon=${horizon}`;
  }

  /**
   * Map API response to ECharts [timestamp, value] pairs.
   * Expected API response shape:
   *   { timestamps: string[], actual: number[], forecast: number[] }
   */
  function toSeriesData(timestamps, values) {
    return timestamps.map((ts, i) => [new Date(ts).getTime(), values[i]]);
  }

  // ── Fetch & render ─────────────────────────────────────────────────────────

  async function fetchAndRender() {
    // Show loading state
    chart.showLoading({ text: "Loading…", maskColor: "rgba(255,255,255,0.6)" });

    try {
      const response = await fetch("http://localhost:8080" + buildURL());

      if (!response.ok) {
        throw new Error(
          `Server returned ${response.status}: ${response.statusText}`,
        );
      }

      const responseData = await response.json();
      const data = responseData.data;
      const timestamps = [],
        actual = [],
        forecast = [];

      if (data.constructor !== Array) {
        throw new Error(`Server returned invalid data`);
      }

      data.map((r) => {
        timestamps.push(r.time);
        actual.push(r.actual);
        forecast.push(r.forecast);
      });

      chart.setOption({
        series: [
          { name: "Actual", data: toSeriesData(timestamps, actual) },
          { name: "Forecast", data: toSeriesData(timestamps, forecast) },
        ],
      });
    } catch (err) {
      console.error("Failed to fetch wind data:", err);
      // Display error inside the chart area so the user is informed
      chart.setOption({
        graphic: [
          {
            type: "text",
            left: "center",
            top: "middle",
            style: {
              text: `⚠ Could not load data\n${err.message}`,
              fontSize: 14,
              fill: "#e53e3e",
              textAlign: "center",
            },
          },
        ],
      });
    } finally {
      chart.hideLoading();
    }
  }

  // ── Event listeners ────────────────────────────────────────────────────────

  // Update the horizon label while dragging; fetch only on release (change)
  slider.addEventListener("input", () => {
    horizonText.innerText = slider.value + " Hr";
  });

  slider.addEventListener("change", fetchAndRender);

  startInput.addEventListener("change", () => {
    // Ensure start never exceeds end
    if (startInput.value > endInput.value) {
      endInput.value = startInput.value;
    }
    fetchAndRender();
  });

  endInput.addEventListener("change", () => {
    // Ensure end never precedes start
    if (endInput.value < startInput.value) {
      startInput.value = endInput.value;
    }
    fetchAndRender();
  });

  // ── Responsiveness ─────────────────────────────────────────────────────────

  window.addEventListener("resize", () => chart.resize());

  // ── Initial load ───────────────────────────────────────────────────────────

  fetchAndRender();
})();
