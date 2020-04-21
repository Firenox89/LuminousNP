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

    httpServer:use('/off', function(req, res)
        level = tonumber(req.query.level)
        print("level " .. level)
        if level == 0 then
            completeLEDBuffer:fill(0,0,0,0)
        elseif level == 1 then
            levelBuffer = ledLevelBuffer[level]
            levelBuffer:fill(0, 0, 0, 0)
            completeLEDBuffer:replace(ledLevelBuffer[level],ledCountPerLevel*level-1)
        elseif level == 2 then
            levelBuffer = ledLevelBuffer[level]
            levelBuffer:fill(0, 0, 0, 0)
            completeLEDBuffer:replace(ledLevelBuffer[level],ledCountPerLevel*level-1)
        elseif level == 3 then
            levelBuffer = ledLevelBuffer[level]
            levelBuffer:fill(0, 0, 0, 0)
            completeLEDBuffer:replace(ledLevelBuffer[level],ledCountPerLevel*level-1)
        elseif level == 4 then
            levelBuffer = ledLevelBuffer[level]
            levelBuffer:fill(0, 0, 0, 0)
            completeLEDBuffer:replace(ledLevelBuffer[level],ledCountPerLevel*level-1)
        end
        ws2812.write(completeLEDBuffer)
        res:sendFile('controls.html')
    end)
    httpServer:use('/update', function(req, res)
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
        end
        res:sendFile('controls.html')
    end)
    httpServer:use('/', function(req, res)
        res:sendFile('controls.html')
    end)
end

setupWebServer()
