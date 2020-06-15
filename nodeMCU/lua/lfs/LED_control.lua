ws2812.init()
--TODO init with config values
local levelCount = 4
local ledCountPerLevel = 47

local effectFileHeaderSize = 12

--setup the led buffers
ledLevelBuffer = {}
ledLevelBuffer[0] = ws2812.newBuffer(levelCount * ledCountPerLevel, 4)
for i = 1, levelCount do
    ledLevelBuffer[i] = ws2812.newBuffer(ledCountPerLevel, 4)
end

local currentEffectTimer

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

function updateLevelBuffer(level, bufferModifyFunc)
    local buffer = ledLevelBuffer[level]
    bufferModifyFunc(buffer)
    if level ~= 0 then
        local offset = ledCountPerLevel * (level - 1) + 1
        ledLevelBuffer[0]:replace(buffer, offset)
    end
    ws2812.write(ledLevelBuffer[0])
end

function on(level)
    stopEffect()
    fillRGBW(level, 0, 0, 0, 255)
end

function off(level)
    stopEffect()
    fillRGBW(level, 0, 0, 0, 0)
end

function fillRGB(level, g, r, b)
    local cg, cr, cb, cw = RGB2RGBW(g, r, b)
    fillRGBW(level, cg, cr, cb, cw)
end

function fillRGBW(level, g, r, b, w)
    stopEffect()
    print(string.format("Fill Level %01d RGBW %03d/%03d/%03d/%03d", level, r, g, b, w))
    updateLevelBuffer(level, function(buffer)
        buffer:fill(g, r, b, w)
    end)
end

function updateRGBW(index, g, r, b, w)
    ledLevelBuffer[0]:set(index, g, r, b, w)
    ws2812.write(ledLevelBuffer[0])
end

function updateBuffer(index, g, r, b, w)
    ledLevelBuffer[0]:set(index, g, r, b, w)
end

function writeBuffer()
    ws2812.write(ledLevelBuffer[0])
end

function bufferShift(timer, delayPerFrame, level)
    local startTime = tmr.now()
    updateLevelBuffer(level, function(buffer)
        buffer:shift(1, ws2812.SHIFT_CIRCULAR)
    end)
    local delta = (tmr.now() - startTime) / 1000
    local delay = delayPerFrame - delta
    --timer overflow
    if (delay < 0) then
        delay = delayPerFrame
    end
    timer:alarm(delay, tmr.ALARM_SINGLE, function()
        bufferShift(timer, delayPerFrame, level)
    end)
end

function playFrames(timer, bytesPerLed, ledCount, frameCount, isRepeating, delayPerFrame, currentFramePos)
    local startTime = tmr.now()
    if (currentFramePos == frameCount) then
        if (isRepeating) then
            file.seek("set", effectFileHeaderSize)
            currentFramePos = 1
        else
            return
        end
    end

    ws2812.write(file.read(bytesPerLed * ledCount))

    local delta = (tmr.now() - startTime) / 1000
    local delay = delayPerFrame - delta
    --timer overflow
    if (delay < 0) then
        delay = delayPerFrame
    end
    timer:alarm(delay, tmr.ALARM_SINGLE, function()
        playFrames(timer, bytesPerLed, ledCount, frameCount, isRepeating, delayPerFrame, currentFramePos + 1)
    end)
end

function interpolateBuffer(factor, buffer1, buffer2, bufferd)
    for i = 1, buffer1:size() do
        --TODO handle RGB buffers
        local g1, r1, b1, w1 = buffer1:get(i)
        local g2, r2, b2, w2 = buffer2:get(i)
        local gd1 = g1 - ((g1 - g2) * factor)
        local rd1 = r1 - ((r1 - r2) * factor)
        local bd1 = b1 - ((b1 - b2) * factor)
        local wd1 = w1 - ((w1 - w2) * factor)
        bufferd:set(i, gd1, rd1, bd1, wd1)
    end
end

function interpolateBufferBatch(factor, buffer1, buffer2, bufferd)
    local dump1 = buffer1:dump()
    local dump2 = buffer2:dump()

    for i = 1, buffer1:size() do
        local offset = ((i - 1) * 4)
        local g1, r1, b1, w1 = string.byte(dump1, offset + 1, offset + 4)
        local g2, r2, b2, w2 = string.byte(dump2, offset + 1, offset + 4)
        --TODO handle RGB buffers
        bufferd:set(i, g1 - ((g1 - g2) * factor), r1 - ((r1 - r2) * factor), b1 - ((b1 - b2) * factor), w1 - ((w1 - w2) * factor))
    end
end

function interpolateFrames(timer, bytesPerLed, frameCount, isRepeating, currentFramePos, buffer1, buffer2, bufferd, interpolationStepCount, interpolationStep)
    local startTime = tmr.now()

    if (currentFramePos == frameCount) then
        if (isRepeating) then
            file.seek("set", effectFileHeaderSize)
            currentFramePos = 1
        else
            return
        end
    end

    if interpolationStep == 0 then
        buffer2:replace(file.read(bytesPerLed * buffer1:size()))
        ws2812.write(buffer1)
    else
        local factor = (1 / interpolationStepCount) * interpolationStep
        interpolateBufferBatch(factor, buffer1, buffer2, bufferd)
        ws2812.write(bufferd)
    end

    local delta = (tmr.now() - startTime) / 1000
    local delay = 70 - delta
    print(delay)
    --timer overflow
    if (delta > delay) then
        delay = 1
    elseif (delay < 0) then
        delay = 70
    end
    local nextStep = (interpolationStep + 1) % interpolationStepCount
    local nextFrame = currentFramePos
    local b1 = buffer1
    local b2 = buffer2
    if (nextStep == 0) then
        nextFrame = nextFrame + 1
        b1 = buffer2
        b2 = buffer1
    end
    timer:alarm(delay, tmr.ALARM_SINGLE, function()
        interpolateFrames(timer, bytesPerLed, frameCount, isRepeating, nextFrame, b1, b2, bufferd, interpolationStepCount, nextStep)
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
        local header = file.read(effectFileHeaderSize)
        --https://nodemcu.readthedocs.io/en/master/modules/struct/#structunpack
        local frameCount, delayPerFrame, bytesPerLed, ledCount, flags = struct.unpack("<hhhhI4", header)

        local repeating = bit.band(flags, bit.bit(0))
        local shiftCircular = bit.band(flags, bit.bit(1))
        local interpolate = bit.band(flags, bit.bit(2))
        local fill = bit.band(flags, bit.bit(3))

        print("Loaded effect file")
        print("File Size " .. size)
        print("Frame Count " .. frameCount)
        print("Delay per frame " .. delayPerFrame)
        print("Bytes Per Led " .. bytesPerLed)
        print("Led count " .. ledCount)

        print("Flags " .. flags)
        print("repeating " .. repeating)
        print("shiftCircular " .. shiftCircular)
        print("interpolate " .. interpolate)
        print("fill " .. fill)

        local level = 0
        local buffer = ledLevelBuffer[level]
        currentEffectTimer = tmr:create()
        if shiftCircular > 0 then
            loadEffectFrameIntoBuffer(bytesPerLed, ledCount, buffer)
            bufferShift(currentEffectTimer, delayPerFrame, level)
            file.close()
        elseif interpolate > 0 then
            local buffer1 = ws2812.newBuffer(ledCount, bytesPerLed)
            local bufferd = ws2812.newBuffer(ledCount, bytesPerLed)
            local buffer2 = ws2812.newBuffer(ledCount, bytesPerLed)
            local interpolationStepCount = math.floor(delayPerFrame / 70)

            interpolateFrames(currentEffectTimer, bytesPerLed, frameCount, repeating, 1, buffer1, buffer2, bufferd, interpolationStepCount, 0)
        else
            playFrames(currentEffectTimer, bytesPerLed, ledCount, frameCount, repeating, delayPerFrame, 1)
        end
    end
end

function stopEffect()
    if (currentEffectTimer ~= nil) then
        currentEffectTimer:unregister()
        currentEffectTimer = nil
    end
    --TODO that might be a bit forceful...
    file.close()
end

return {
    on = on,
    off = off,
    fill = fillRGB,
    updateRGBW = updateRGBW,
    playEffect = playEffect,
    stopEffect = stopEffect
}
