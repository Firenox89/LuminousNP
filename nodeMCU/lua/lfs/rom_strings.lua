
local preload="\
WiFi connection to AP("," of ","(",")",") established!",") has failed!","<form action=\"/effect\" method=\"get\"><label for=\"effect\">Choose a effect:</label><select id=\"effect\" name=\"effect\"><option value=\"rainbow\">Rainbow</option><option value=\"rainbowroad\">RainbowRoad</option></select><input type=\"hidden\" id=\"level\" name=\"level\" value=\"%01d\"><input type=\"submit\" value=\"Apply\"></form>","<form action=\"/fill\" method=\"get\"><label for=\"color\">Color:</label><input type=\"color\" id=\"color\" name=\"color\" value=\"#%02X%02X%02X\"><br><br><label for=\"white\">White:</label> <input type=\"range\" id=\"white\" name=\"white\" value=\"%01d\" min=\"0\" max=\"255\"><br><br><input type=\"hidden\" id=\"level\" name=\"level\" value=\"%01d\"><input type=\"submit\" value=\"Update\"></form>","<h2>Level %01d</h2><a href=\"off?level=%01d\"><button>Off</button></a><br><br>","ALARM_SINGLE","ASSOC_LEAVE","Aborting connection to AP!","Check for update","Connection to AP(","Disconnect reason: ","GET","GET /debug HTTP/1.1\r\
Host: 192.168.2.101\r\
User-Agent: curl/7.69.1\r\
Accept: */*\r\
\r\
","IP","Lua 5.1","Retrying connection...(attempt ","Running","Start application","Startup will resume momentarily, you have 3 seconds to abort.","Waiting for IP address...","Waiting...","Wifi connection is ready! IP address is: ","_VERSION","config","crypto.hash","disconnect","disconnect_ct","eventmon","init.lua","init.lua deleted or renamed","kv","reason","sta","startup","tls.socket","websocket.client","wifi_connect_event","wifi_disconnect_event","wifi_got_ip_event","ws2812.buffer"