# weather_demo
This is demo project to get weather data by golang.

Run this project. First build executable file:

    go build

Then you can execute binary file:

    weather_demo.exe (windows)

Server would listen port on 8000 and provide three API:

  - POST http://localhost:8000/v1/login
  - GET  http://localhost:8000/v1/weather
  - POST http://localhost:8000/v1/logout

More details about API request example please refer to test/ folder.
There is a postman file.
