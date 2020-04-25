ws2812.init()
local levelCount = 4
local ledCountPerLevel = 47
ledLevelBuffer = {}
local effectCoroutineLevels = {}

ledLevelBuffer[0] = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)

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
        while ( true ) do
            updateAndFlashLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    hue = 0
                    saturation = 255
                    brightness = 255
                    bufferFilled = false
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

function rainbowSnake(level)
    co = coroutine.create(function ()
        bufferFilled = false
        while ( true ) do
            updateAndFlashLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    hue = 0
                    saturation = 255
                    brightness = 255
                    inc = 10
                    for i = 1, 36 do
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

function runningLight(level)
    co = coroutine.create(function ()
        bufferFilled = false
        while ( true ) do
            updateAndFlashLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    g, r, b, w = buffer:get(1)
                    for i = 0, 5 do
                        brightness = 5 + (i * 50)
                        g, r, b, w = color_utils.hsv2grbw(0, 255, brightness)
                        buffer:set(i+1, g, r, b, w)
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

function updateAndFlashLevelBuffer(level, bufferModifyFunc)
    local buffer
    buffer = ledLevelBuffer[level]
    bufferModifyFunc(buffer)
    if level ~= 0 then
        offset = ledCountPerLevel*(level-1) + 1
        ledLevelBuffer[0]:replace(buffer, offset)
    end
    ws2812.write(ledLevelBuffer[0])
    ws2812.write(ledLevelBuffer[0])
end

function setupWebServer()
    print("setup webserver")
    dofile('httpServer.lua')
    httpServer:listen(80)

    httpServer:use('/OTA', function(req, res)
        file.open("OTA.update", "w")
        file.close()
        node.restart()
        res:send('Fetching')
    end)

    httpServer:use('/telnet', function(req, res)
        LFS.telnet()
        res:send('Telnet started')
    end)

    httpServer:use('/debugStrings', function(req, res)
        do
            local a=debug.getstrings'RAM'
            for i =1, #a do a[i] = ('%q'):format(a[i]) end
            print ('local preload='..table.concat(a,','))
        end
        res:send('Done')
    end)

    httpServer:use('/on', function(req, res)
        level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end

        if level == 0 then
            for i = 0, levelCount do
                effectCoroutineLevels[i] = nil
            end
        else
            effectCoroutineLevels[level] = nil
        end
        updateAndFlashLevelBuffer(level, function (buffer)
            buffer:fill(0, 0, 0, 255)
        end)
        res:send(200)
    end)

    httpServer:use('/off', function(req, res)
        level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end

        if level == 0 then
            for i = 0, levelCount do
                effectCoroutineLevels[i] = nil
            end
        else
            effectCoroutineLevels[level] = nil
        end
        updateAndFlashLevelBuffer(level, function (buffer)
            buffer:fill(0, 0, 0, 0)
        end)
        res:send(200)
    end)

    httpServer:use('/fill', function(req, res)
        level = tonumber(req.query.level)
        rawColor = req.query.color
        r = tonumber(rawColor:sub(1, 2),16)
        g = tonumber(rawColor:sub(3, 4),16)
        b = tonumber(rawColor:sub(5, 6),16)
        w = tonumber(req.query.white)
        if r and r >= 0 and r <= 255
            and b and b >= 0 and b <= 255
            and g and g >= 0 and g <= 255
            and w and w >= 0 and w <= 255
            and level and level >= 0 and level <= 4
            then
            updateAndFlashLevelBuffer(level, function (buffer)
                print(string.format("Fill Level %01d RGBW %03d/%03d/%03d/%03d", level, r, g, b, w))
                buffer:fill(g, r, b, w)
            end)
        end
        res:send(200)
    end)

    httpServer:use('/effects', function(req, res)
        effectJson = "[\"Rainbow\", \"Rainbow Road\", \"Rainbow Snake\", \"Running Light\"]"
        res:send(effectJson)
    end)

    httpServer:use('/effect', function(req, res)
        level = tonumber(req.query.level)
        effect = req.query.effect
        rawColor = req.query.color
        r = tonumber(rawColor:sub(1, 2),16)
        g = tonumber(rawColor:sub(3, 4),16)
        b = tonumber(rawColor:sub(5, 6),16)
        --hue, sat, bri = color_utils.grb2hsv(50,15,200)
        if level == nil then
            level = 0
        end
        if level == 0 then
            for i = 0, levelCount do
                effectCoroutineLevels[i] = nil
            end
        end
        print("set " .. effect or "no effect" .. " for " .. level)
        if effect == "Rainbow" then
            effectCoroutineLevels[level] = rainbow(level)
        elseif effect == "Rainbow Road" then
            effectCoroutineLevels[level] = rainbowRoad(level)
        elseif effect == "Rainbow Snake" then
            effectCoroutineLevels[level] = rainbowSnake(level)
        elseif effect == "Running Light" then
            effectCoroutineLevels[level] = runningLight(level)
        end
        res:send(200)
    end)
end

function effectLoop()
    for i = 0, levelCount do
        effectCoroutine = effectCoroutineLevels[i]
        if effectCoroutine ~= nil then
            coroutine.resume(effectCoroutine)
        end
    end
end

if file.exists("OTA.update") then
    print("OTA file exists")
    file.remove("OTA.update")
    LFS.HTTP_OTA("sirmixalot", "/", "lfs.img")
else
    tmr.create():alarm(50, tmr.ALARM_AUTO, effectLoop)

    setupWebServer()
end

