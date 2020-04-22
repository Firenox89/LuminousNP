ws2812.init()
levelCount = 4
ledCountPerLevel = 47
ledLevelBuffer = {}
effectCoroutineLevels = {}

local effectCoroutineAll
completeLEDBuffer = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)

webHeader = "<html> <body> <a href=\"off?level=0\"><button>All Off</button></a>"
formPart1 = "<h2>Level " --level
formPart2 = "</h2><a href=\"off?level=" --level
formPart3 = "\"><button>Off</button></a><br><br> <form action=\"/fill\" method=\"get\"> <label for=\"color\">Color:</label>"
formPart4 = "<input type=\"color\" id=\"color\" name=\"color\" value=\"#" --color
formPart5 = "\"><br><br> <label for=\"white\">White:</label> <input type=\"range\" id=\"white\" name=\"white\" value=\"" --white
formPart6 = "\" min=\"0\" max=\"255\"><br><br> <input type=\"hidden\" id=\"level\" name=\"level\" value=\"" --level
formPart7 = "\"><input type=\"submit\" value=\"Update\"></form>"
webFooter = "</body> </html>"

for i = 1, levelCount do
    ledLevelBuffer[i] = ws2812.newBuffer(ledCountPerLevel,4)
    effectCoroutineLevels[i] = nil
end

function rainbow(level)
    co = coroutine.create(function ()
        hue = 0
        saturation = 255
        brightness = 255
        while ( true ) do
            hue = (hue + 7) % 360
            g, r, b, w = color_utils.hsv2grbw(hue, saturation, brightness)
            updateAndFlashLevelBuffer(level, function (buffer)
                buffer:fill(g, r, b, w)
            end)
            coroutine.yield()
        end
    end)
    return co
end

function colorwipe()
    co = coroutine.create(function ()
        g = 0
        r = 0
        b = 0
        w = 0
        while ( true ) do
            g, r, b, w = color_utils.hsv2grbw(hue, saturation, brightness)
            completeLEDBuffer:fill(g, r, b, w)
            ws2812.write(completeLEDBuffer)
            coroutine.yield()
        end
    end)
    return co
end

function updateAndFlashLevelBuffer(level, bufferModifyFunc)
    local buffer
    if level == 0 then
        buffer = completeLEDBuffer
    else
        buffer = ledLevelBuffer[level]
    end
    bufferModifyFunc(buffer)
    if level ~= 0 then
        offset = ledCountPerLevel*(level-1) + 1
        completeLEDBuffer:replace(buffer, offset)
    end
    ws2812.write(completeLEDBuffer)
end

function sendPage(res)
    white = 255
    colorString = "ff0000"
    buf = webHeader
    for level = 1, levelCount do
        buf = buf .. formPart1 .. level .. formPart2 .. level .. formPart3 .. formPart4 .. colorString .. formPart5 .. white .. formPart6 .. level .. formPart7
    end
    buf = buf .. webFooter

    res:send(buf)
end

function setupWebServer()
    print("setup webserver")
    dofile('httpServer.lua')
    httpServer:listen(80)

    httpServer:use('/', function(req, res)
        sendPage(res)
    end)

    httpServer:use('/OTA', function(req, res)
        LFS.HTTP_OTA("sirmixalot", "/", "lfs.img")
        res:send('Fetching')
    end)

    httpServer:use('/off', function(req, res)
        level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end
        effectCoroutineAll = nil

        for i = 1, levelCount do
            effectCoroutineLevels[i] = nil
        end
        updateAndFlashLevelBuffer(level, function (buffer)
            buffer:fill(0, 0, 0, 0)
        end)
        sendPage(res)
    end)
    httpServer:use('/fill', function(req, res)
        rawColor = req.query.color
        level = tonumber(req.query.level)
        r = tonumber(rawColor:sub(2, 3),16)
        g = tonumber(rawColor:sub(4, 5),16)
        b = tonumber(rawColor:sub(6, 7),16)
        w = tonumber(req.query.white)
        print(string.format("Level %01d RGBW %03d/%03d/%03d/%03d", level, r, g, b, w))
        if r and r >= 0 and r <= 255
            and b and b >= 0 and b <= 255
            and g and g >= 0 and g <= 255
            and w and w >= 0 and w <= 255
            and level and level > 0 and level <= 4
            then
            updateAndFlashLevelBuffer(level, function (buffer)
                buffer:fill(g, r, b, w)
            end)
            print("color updated")
        end
        sendPage(res)
    end)
    httpServer:use('/rainbow', function(req, res)
        level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end
        print("set rainbow for " .. level)
        if level == 0 then
            effectCoroutineAll =  rainbow(level)
        else
            effectCoroutineLevels[level] = rainbow(level)
        end
        sendPage(res)
    end)
end

function effectLoop()
    if effectCoroutineAll ~= nil then
        coroutine.resume(effectCoroutineAll)
    end

    for i = 1, levelCount do
        effectCoroutine = effectCoroutineLevels[i]
        if effectCoroutine ~= nil then
            coroutine.resume(effectCoroutine)
        end
    end
end

tmr.create():alarm(50, tmr.ALARM_AUTO, effectLoop)

setupWebServer()

