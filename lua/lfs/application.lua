function startLEDs()
    print("start LEDS")
    ws2812.init()
    ws2812.write(string.char(255, 0, 0, 0, 255, 0, 0, 0))
end

function setupWebServer()
    print("setup webserver")
    dofile('httpServer.lua')
    httpServer:listen(80)

    httpServer:use('/OTA', function(req, res)
        LFS.HTTP_OTA("sirmixalot", "/", "lfs.img")
        res:send('Fetching')
    end)

    levelCount = 4
    ledCountPerLevel = 47
    ledLevelBuffer = {}

    completeLEDBuffer = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)

    for i = 1, levelCount do
        ledLevelBuffer[i] = ws2812.newBuffer(ledCountPerLevel,4)
    end

    ws2812.init()

    httpServer:use('/test', function(req, res)
        print("test")
        res:send('test')
    end)
    httpServer:use('/off', function(req, res)
        buffer = ws2812.newBuffer(200, 4)
        buffer:fill(0,0,0,0)
        ws2812.write(buffer)
        res:send('Off')
    end)
    httpServer:use('/fill', function(req, res)
        level = tonumber(req.query.level)
        r = tonumber(req.query.r)
        g = tonumber(req.query.g)
        b = tonumber(req.query.b)
        w = tonumber(req.query.w)
        if r and r >= 0 and r <= 255
            and b and b >= 0 and b <= 255
            and g and g >= 0 and g <= 255
            and w and w >= 0 and w <= 255
            and level and level >= 0 and level <= 3
            then
            levelBuffer = ledLevelBuffer[level+1]
            levelBuffer:fill(r, g, b, w)
            completeLEDBuffer:replace(ledLevelBuffer[level+1],ledCountPerLevel*level+1)
            ws2812.write(completeLEDBuffer)
            res:send("ok")
        end
    end)
    httpServer:use('/white', function(req, res)
        completeLEDBuffer:fill(255,255,255,255)
        ws2812.write(completeLEDBuffer)
        res:send(completeLEDBuffer:dump())
    end)
    httpServer:use('/red', function(req, res)
        completeLEDBuffer:fill(0,255,0,111)
        ws2812.write(completeLEDBuffer)
        res:send(completeLEDBuffer:dump())
    end)
    httpServer:use('/blau1', function(req, res)
        ledLevelBuffer[1]:fill(0,255,0,0)
        completeLEDBuffer:replace(ledLevelBuffer[1],ledCountPerLevel)
        ws2812.write(completeLEDBuffer)
        res:send('red')
    end)
end

-- a simple telnet server
s=net.createServer(net.TCP)
s:listen(2323,function(c)
   con_std = c
   function s_output(str)
      if(con_std~=nil)
         then con_std:send(str)
      end
   end
   node.output(s_output, 0)   -- re-direct output to function s_ouput.
   c:on("receive",function(c,l)
      node.input(l)           -- works like pcall(loadstring(l)) but support multiple separate line
   end)
   c:on("disconnection",function(c)
      con_std = nil
      node.output(nil)        -- un-regist the redirect output function, output goes to serial
   end)
end)

setupWebServer()
