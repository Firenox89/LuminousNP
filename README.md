# LuminousNP Overview

## Needed Hardware
- RaspberryPi of any sort
- NodeMCU (ESP8266)
- LEDs Strips (WS2812B or similar or SK6812)
- The stuff to power and connect everything

## Software contained in this Project

### NodeMCU

- The NodeMCU Firmware with the needed modules
- The lua code to connect to the local WiFi, provide the control API and the LED control

### RaspberryPI

- The backend to host the Web Interface and connect the NodeMCUs via WiFi
- The WebInterface as a ReactJS Project
- The code to power the Touchscreen, to controll the setup without a browser
