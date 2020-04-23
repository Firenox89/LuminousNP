ws2812.init()
local levelCount = 4
local ledCountPerLevel = 47
local ledLevelBuffer = {}
local effectCoroutineLevels = {}

local effectCoroutineAll
local completeLEDBuffer = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)

local webHeader = "<html> <body>"
local offButton = "<h2>Level %01d</h2>" ..
"<a href=\"off?level=%01d\"><button>Off</button></a><br><br>"
local effectForm = "<form action=\"/effect\" method=\"get\">" ..
"<label for=\"effect\">Choose a effect:</label>" ..
"<select id=\"effect\" name=\"effect\">" ..
"<option value=\"rainbow\">Rainbow</option>" ..
"<option value=\"rainbowroad\">RainbowRoad</option>" ..
"</select>" ..
"<input type=\"hidden\" id=\"level\" name=\"level\" value=\"%01d\">" ..
"<input type=\"submit\" value=\"Apply\">" ..
"</form>"
local colorForm = "<form action=\"/fill\" method=\"get\">" ..
"<label for=\"color\">Color:</label>" ..
"<input type=\"color\" id=\"color\" name=\"color\" value=\"#%02X%02X%02X\"><br><br>" ..
"<label for=\"white\">White:</label> <input type=\"range\" id=\"white\" name=\"white\" value=\"%01d\" min=\"0\" max=\"255\"><br><br>" ..
"<input type=\"hidden\" id=\"level\" name=\"level\" value=\"%01d\">" ..
"<input type=\"submit\" value=\"Update\">" ..
"</form>"
local webFooter = "</body> </html>"

for i = 1, levelCount do
    ledLevelBuffer[i] = ws2812.newBuffer(ledCountPerLevel,4)
    effectCoroutineLevels[i] = nil
end

function sendPage(res)
    buf = webHeader
    buf = buf .. string.format(offButton, 0, 0)
    buf = buf .. string.format(effectForm, 0)
    g, r, b, w = completeLEDBuffer:get(1)
    buf = buf .. string.format(colorForm, r, g, b, w, 0)
    for level = levelCount, 1, -1 do
        g, r, b, w = ledLevelBuffer[level]:get(1)
        buf = buf .. string.format(offButton, level, level)
        buf = buf .. string.format(effectForm, level)
        buf = buf .. string.format(colorForm, r, g, b, w, level)
    end
    buf = buf .. webFooter

    res:send(buf)
end

function rainbow(level)
    co = coroutine.create(function ()
        hue = 0
        saturation = 255
        brightness = 255
        while ( true ) do
            hue = (hue + 1) % 360
            g, r, b, w = color_utils.hsv2grbw(hue, saturation, brightness)
            updateAndFlashLevelBuffer(level, function (buffer)
                buffer:fill(g, r, b, w)
            end)
            coroutine.yield()
        end
    end)
    return co
end

function rainbowRoad(level)
    co = coroutine.create(function ()
        hue = 0
        saturation = 255
        brightness = 255
        bufferFilled = false
        while ( true ) do
            updateAndFlashLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    inc = 360 / buffer:size()
                    for i = 1, buffer:size() do
                        hue = (i * inc) % 360
                        g, r, b, w = color_utils.hsv2grbw(hue, saturation, brightness)
                        buffer:set(i, g, r, b, w)
                    end
                    bufferFilled = true
                    print("buffer filled")
                end
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

    httpServer:use('/debugStrings', function(req, res)
        do
            local a=debug.getstrings'RAM'
            for i =1, #a do a[i] = ('%q'):format(a[i]) end
            print ('local preload='..table.concat(a,','))
        end
        res:send('Done')
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
        print(string.format("Fill Level %01d RGBW %03d/%03d/%03d/%03d", level, r, g, b, w))
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

    httpServer:use('/effect', function(req, res)
        level = tonumber(req.query.level)
        effect = req.query.effect
        print(effect)
        if level == nil then
            level = 0
        end
        local effectFunc
        if effect == "rainbow" then
            effectFunc = rainbow
        elseif effect == "rainbowroad" then
            effectFunc = rainbowRoad
        end
        if effectFunc ~= nil then
            print("set " .. effect .. " for " .. level)
            if level == 0 then
                effectCoroutineAll = effectFunc(level)
            else
                effectCoroutineLevels[level] = effectFunc(level)
            end
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

