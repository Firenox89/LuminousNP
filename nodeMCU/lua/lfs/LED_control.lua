ws2812.init()
local levelCount = 4
local ledCountPerLevel = 47

local effectCoroutineLevels = {}

--setup the led buffers
ledLevelBuffer = {}
ledLevelBuffer[0] = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)
for i = 1, levelCount do
    ledLevelBuffer[i] = ws2812.newBuffer(ledCountPerLevel, 4)
    effectCoroutineLevels[i] = nil
end

function RGB2RGBW(r, g, b)
    local cg, cr, cb, cw
    if (r ~= g or r ~= b) then
        local hue, saturation, brightness
        hue, saturation, brightness = color_utils.grb2hsv(g, r, b)
        cg, cr, cb, cw = color_utils.hsv2grbw(hue, saturation, brightness)
    else
        cr = r
        cg = r
        cb = r
        cw = r
    end
    return cg, cr, cb, cw
end

function rainbow(level)
    return coroutine.create(function()
        local hue = 0
        local saturation = 255
        local brightness = 255
        while (true) do
            hue = (hue + 1) % 360
            g, r, b, w = color_utils.hsv2grbw(hue, saturation, brightness)
            updateLevelBuffer(level, function(buffer)
                buffer:fill(g, r, b, w)
            end)
            coroutine.yield()
        end
    end)
end

function rainbowRoad(level)
    return coroutine.create(function()
        local bufferFilled = false
        while (true) do
            updateLevelBuffer(level, function(buffer)
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
end

function rainbowSnake(level)
    return coroutine.create(function()
        local bufferFilled = false
        while (true) do
            updateLevelBuffer(level, function(buffer)
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
end

function runningLight(level, hue)
    return coroutine.create(function()
        local bufferFilled = false
        while (true) do
            updateLevelBuffer(level, function(buffer)
                if (bufferFilled) then
                    buffer:shift(1, ws2812.SHIFT_CIRCULAR)
                else
                    for i = 0, 5 do
                        local brightness = 5 + (i * 50)
                        g, r, b, w = color_utils.hsv2grbw(hue, 255, brightness)
                        buffer:set(i + 1, g, r, b, w)
                    end
                    bufferFilled = true
                    print("buffer filled")
                end
            end)
            coroutine.yield()
        end
    end)
end

function pulsingLight(level, hue)
    return coroutine.create(function()
        local brightnessInc = 3
        local brightness = 0
        while (true) do
            if (brightness + brightnessInc > 255 or brightness + brightnessInc < 0) then
                brightnessInc = -brightnessInc
            end
            brightness = brightness + brightnessInc
            updateLevelBuffer(level, function(buffer)
                g, r, b, w = color_utils.hsv2grbw(hue, 255, brightness)
                buffer:fill(g, r, b, w)
            end)
            coroutine.yield()
        end
    end)
end

function updateLevelBuffer(level, bufferModifyFunc)
    local buffer = ledLevelBuffer[level]
    bufferModifyFunc(buffer)
    if level ~= 0 then
        local offset = ledCountPerLevel * (level - 1) + 1
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
    fillRGBW(level, 0, 0, 0, 255)
end

function off(level)
    fillRGBW(level, 0, 0, 0, 0)
end

function fillRGB(level, g, r, b)
    local cg, cr, cb, cw = RGB2RGBW(g, r, b)
    fillRGBW(level, cg, cr, cb, cw)
end

function fillRGBW(level, g, r, b, w)
    print(string.format("Fill Level %01d RGBW %03d/%03d/%03d/%03d", level, r, g, b, w))
    if level == 0 then
        for i = 0, levelCount do
            effectCoroutineLevels[i] = nil
        end
    else
        effectCoroutineLevels[level] = nil
    end
    updateLevelBuffer(level, function(buffer)
        buffer:fill(g, r, b, w)
    end)
end

function updateRGBW(index, g, r, b, w)
    for i = 0, levelCount do
        effectCoroutineLevels[i] = nil
    end
    ledLevelBuffer[0]:set(index, g, r, b, w)
    ws2812.write(ledLevelBuffer[0])
end

function updateBuffer(index, g, r, b, w)
    ledLevelBuffer[0]:set(index, g, r, b, w)
end

function writeBuffer()
    ws2812.write(ledLevelBuffer[0])
end

function setEffect(level, effectCoroutine)
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

function bufferShift(level)
    return coroutine.create(function()
        while (true) do
            updateLevelBuffer(level, function(buffer)
                buffer:shift(1, ws2812.SHIFT_CIRCULAR)
            end)
            coroutine.yield()
        end
    end)
end

function playFrames(bytesPerLed, ledCount, level)
    return coroutine.create(function()
        while (true) do
            updateLevelBuffer(level, function(buffer)
                loadEffectFrameIntoBuffer(bytesPerLed, ledCount, buffer)
            end)
            coroutine.yield()
        end
    end)
end

function loadEffectFrameIntoBuffer(bytesPerLed, ledCount, ledbuffer)
    local frameValues = file.read(bytesPerLed * ledCount)
    for j = 1, ledCount do
        local offset = ((j - 1) * bytesPerLed)
        local r = string.byte(frameValues, offset + 1)
        local g = string.byte(frameValues, offset + 2)
        local b = string.byte(frameValues, offset + 3)
        local w = string.byte(frameValues, offset + 4)

        ledbuffer:set(j, g, r, b, w)
    end
end

function playEffect()
    -- print the first 5 bytes of 'init.lua'
    local fd = file.open("current.effect", "r")
    if fd then
        local size = file.stat("current.effect").size
        local header = file.read(12)
        --https://nodemcu.readthedocs.io/en/master/modules/struct/#structunpack
        local schemaVersion, delayPerFrame, bytesPerLed, ledCount, flags = struct.unpack("<hhhhI4", header)

        local isRepeating = bit.band(flags, bit.bit(0))
        local isShiftCircular = bit.band(flags, bit.bit(1))

        local frameCount = (size - 12) / (bytesPerLed * ledCount)

        print("Loaded effect file")
        print("Schema " .. schemaVersion)
        print("Frame Count " .. frameCount)
        print("Delay per frame " .. delayPerFrame)
        print("Bytes Per Led " .. bytesPerLed)
        print("Led count " .. ledCount)

        print("Flags " .. flags)
        print("IsRepeating " .. isRepeating)
        print("IsShiftCircular " .. isShiftCircular)

        local level = 0
        local buffer = ledLevelBuffer[level]
        if isShiftCircular > 0 then
            loadEffectFrameIntoBuffer(bytesPerLed, ledCount, buffer)
            setEffect(level, bufferShift(level))
            file.close()
        else
            setEffect(level, playFrames(bytesPerLed, ledCount, level))
        end
    end
end

function stopEffect()
    effectCoroutineLevels[0] = nil
    file.close()
end

return {
    init = init,
    on = on,
    off = off,
    fill = fillRGB,
    updateRGBW = updateRGBW,
    setRainbow = setRainbow,
    setRainbowRoad = setRainbowRoad,
    setRainbowSnake = setRainbowSnake,
    setRunningLight = setRunningLight,
    setPulsingLight = setPulsingLight,
    playEffect = playEffect,
    stopEffect = stopEffect
}
