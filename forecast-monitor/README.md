# Forecast Monitoring App

This app plots 2 lines - one for the actual wind power and other for the forecasted wind power. User can select the timestamp range and the forecast horizon to monitor the forecasting results against actual data.

Few terminologies:

1. Target time: The time when the event happens. The API treats "start time" as the target time. The target time specifies the 30 min span of forecast data. Example, startTime = 2024-01-05 18:00 means the wind generation between 18:00 - 18:30.
2. Publish time: When the forecast API is called, it returns the publish time. It is the time when the system created the forecast.
3. Forecast horizon: forecast horizon = startTime - publishTime. Example: publishTime = 10:00, startTime = 18:00, horizon = 8 hours. This is an 8-hour ahead forecast.

Every data point in the forecast line (green line) represents the latest forecast that was created at least X hours before the target time. Here, X is the forecast horizon.

**Technologies**

- Frontend is written is HTML, CSS and JS.
- Charts are created using Apache Echarts.

**AI Usage**

1. Generate frontend code and styles.
2. Research terminologies in for the wind data.
3. Discuss the whether to divide the architecture into frontend and backend.
4. List the important tests that need to be done for the API.
