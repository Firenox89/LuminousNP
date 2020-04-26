ws2812.init()
local levelCount = 4
local ledCountPerLevel = 47

local effectCoroutineLevels = {}

--setup the led buffers
ledLevelBuffer = {}
ledLevelBuffer[0] = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)
for i = 1, levelCount do
    ledLevelBuffer[i] = ws2812.newBuffer(ledCountPerLevel,4)
    effectCoroutineLevels[i] = nil
end

function rainbow(level)
    co = coroutine.create(function ()
        local hue = 0
        local saturation = 255
        local brightness = 255
        while ( true ) do
            hue = (hue + 1) % 360
            g, r, b, w = color_utils.hsv2grbw(hue, saturation, brightness)
            updateLevelBuffer(level, function (buffer)
                buffer:fill(g, r, b, w)
            end)
            coroutine.yield()
        end
    end)
    return co
end

function rainbowRoad(level)
    co = coroutine.create(function ()
        local bufferFilled = false
        while ( true ) do
            updateLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    local hue = 0
                    local saturation = 255
                    local brightness = 255
                    local inc = 360 / buffer:size()
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
        local bufferFilled = false
        while ( true ) do
            updateLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    local hue = 0
                    local saturation = 255
                    local brightness = 255
                    local inc = 10
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

function runningLight(level, hue)
    co = coroutine.create(function ()
        local bufferFilled = false
        while ( true ) do
            updateLevelBuffer(level, function (buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    for i = 0, 5 do
                        local brightness = 5 + (i * 50)
                        g, r, b, w = color_utils.hsv2grbw(hue, 255, brightness)
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

function pulsingLight(level, hue)
    co = coroutine.create(function ()
        local brightnessInc = 3
        local brightness = 0
        while ( true ) do
            if (brightness + brightnessInc > 255 or brightness + brightnessInc < 0) then
                brightnessInc = -brightnessInc
            end
            brightness = brightness + brightnessInc
            updateLevelBuffer(level, function (buffer)
                g, r, b, w = color_utils.hsv2grbw(hue, 255, brightness)
                buffer:fill(g, r, b, w)
            end)
            coroutine.yield()
        end
    end)
    return co
end

function updateLevelBuffer(level, bufferModifyFunc)
    local buffer = ledLevelBuffer[level]
    bufferModifyFunc(buffer)
    if level ~= 0 then
        local offset = ledCountPerLevel*(level-1) + 1
        ledLevelBuffer[0]:replace(buffer, offset)
    end
end

function effectLoop()
    for i = 0, levelCount do
        local effectCoroutine = effectCoroutineLevels[i]
        if effectCoroutine ~= nil then
            coroutine.resume(effectCoroutine)
        end
    end
    ws2812.write(ledLevelBuffer[0])
end

function init()
    tmr.create():alarm(50, tmr.ALARM_AUTO, effectLoop)
end

function on(level)
    fill(level, 0, 0, 0, 255)
end

function off(level)
    fill(level, 0, 0, 0, 0)
end

function fill(level, g, r, b, w)
    print(string.format("Fill Level %01d RGBW %03d/%03d/%03d/%03d", level, r, g, b, w))
    if level == 0 then
        for i = 0, levelCount do
            effectCoroutineLevels[i] = nil
        end
    else
        effectCoroutineLevels[level] = nil
    end
    updateLevelBuffer(level, function (buffer)
        buffer:fill(g, r, b, w)
    end)
end

function setEffect(effectCoroutine)
    if level == 0 then
        for i = 0, levelCount do
            effectCoroutineLevels[i] = nil
        end
    end
    effectCoroutineLevels[level] = effectCoroutine
end

function setRainbow(level)
    setEffect(rainbow(level))
end

function setRainbowRoad(level)
    setEffect(rainbowRoad(level))
end

function setRainbowSnake(level)
    setEffect(rainbowSnake(level))
end

function setRunningLight(level, hue)
    setEffect(runningLight(level, hue))
end

function setPulsingLight(level, hue)
    setEffect(pulsingLight(level, hue))
end

return {
    init = init,
    on = on,
    off = off,
    fill = fill,
    setRainbow = setRainbow,
    setRainbowRoad = setRainbowRoad,
    setRainbowSnake = setRainbowSnake,
    setRunningLight = setRunningLight,
    setPulsingLight = setPulsingLight
}
