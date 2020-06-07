LEDs = require 'LED_control'

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

    httpServer:use('/test', function(req, res)
        res:send('Test')
    end)

    httpServer:use('/debugStrings', function(req, res)
        do
            local a = debug.getstrings 'RAM'
            for i = 1, #a do
                a[i] = ('%q'):format(a[i])
            end
            print('local preload=' .. table.concat(a, ','))
        end
        res:send('Done')
    end)

    httpServer:use('/on', function(req, res)
        level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end
        LEDs.on(level)

        res:send(200)
    end)

    httpServer:use('/off', function(req, res)
        level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end

        LEDs.off(level)
        res:send(200)
    end)

    httpServer:use('/fill', function(req, res)
        level = tonumber(req.query.level)
        rawColor = req.query.color
        r = tonumber(rawColor:sub(1, 2), 16)
        g = tonumber(rawColor:sub(3, 4), 16)
        b = tonumber(rawColor:sub(5, 6), 16)
        if r and r >= 0 and r <= 255
                and b and b >= 0 and b <= 255
                and g and g >= 0 and g <= 255
                and level and level >= 0 and level <= 4
        then
            LEDs.fill(level, r, g, b)
        end
        res:send(200)
    end)

    httpServer:use('/playEffect', function(req, res)
        playEffect()
        res:send(200)
    end)

    httpServer:use('/effects', function(req, res)
        effectJson = "[\"Rainbow\", \"Rainbow Road\", \"Rainbow Snake\", \"Running Light\", \"Pulsing Light\"]"
        res:send(effectJson)
    end)

    httpServer:use('/effect', function(req, res)
        level = tonumber(req.query.level)
        effect = req.query.effect
        rawColor = req.query.color
        r = tonumber(rawColor:sub(1, 2), 16)
        g = tonumber(rawColor:sub(3, 4), 16)
        b = tonumber(rawColor:sub(5, 6), 16)
        hue, sat, bri = color_utils.grb2hsv(g, r, b)
        if level == nil then
            level = 0
        end
        print("set " .. effect or "no effect" .. " for " .. level)
        if effect == "Rainbow" then
            LEDs.setRainbow(level)
        elseif effect == "Rainbow Road" then
            LEDs.setRainbowRoad(level)
        elseif effect == "Rainbow Snake" then
            LEDs.setRainbowSnake(level)
        elseif effect == "Running Light" then
            LEDs.setRunningLight(level, hue)
        elseif effect == "Pulsing Light" then
            LEDs.setPulsingLight(level, hue)
        end
        res:send(200)
    end)

    httpServer:use('/presets', function(req, res)
        effectJson = "[\"Rainbow\", \"Rainbow Road\", \"Rainbow Snake\", \"Running Light\", \"Pulsing Light\"]"
        res:send(effectJson)
    end)

    httpServer:use('/loadPreset', function(req, res)
        preset = req.query.preset
        print("Load " .. preset)
    end)

    httpServer:use('/savePreset', function(req, res)
        preset = req.query.preset
        print("Save " .. preset)
    end)

    httpServer:use('/startStream', function(req, res)
        openUDPSocket()
    end)

    registerAtController()
end

local isConnectedToController = false
local reconnectionDelay = 5000

function registerAtController()
    print("Try to register at controller")
    if not isConnectedToController then
        local socket = net.createConnection()
        socket:on("connection", function(sck, c)
            print("controller connected")
            socket:send('{"id": 2, "ledCount": 188, "bytesPerLed": 4}')
            isConnectedToController = true
        end)
        socket:on("disconnection", function(sck, c)
            print("controller disconnected")
            isConnectedToController = false
            if not tmr.create():alarm(reconnectionDelay, tmr.ALARM_SINGLE, function()
                registerAtController()
            end)
            then
                print("Failed to start reconnection timer.")
            end
        end)
        socket:connect(4488, "nodemcu-controller")
    else
        print("Already connected")
    end
end

function playEffect()
    -- print the first 5 bytes of 'init.lua'
    local fd = file.open("current.effect", "r")
    if fd then
        local size = file.stat("current.effect").size
        local header = fd:read(12)
        --https://nodemcu.readthedocs.io/en/master/modules/struct/#structunpack
        local schemaVersion, delayPerFrame, bytesPerLed, ledCount = struct.unpack("<hhhhI4", header)

        local frameCount = (size - 12) / (bytesPerLed * ledCount)

        print("schema " .. schemaVersion .. " frameCounts " .. frameCount .. " deley per frame " .. delayPerFrame .. " bytesPerLed " .. bytesPerLed .. " led count " .. ledCount)

        for i = 1, frameCount do
            print(string.format("Process Frame %d", i))

            local frameValues = fd:read(bytesPerLed*ledCount)
            for j = 1, ledCount do
                local offset = ((j-1) * bytesPerLed)
                local r = string.byte(frameValues, offset + 1)
                local g = string.byte(frameValues, offset + 2)
                local b = string.byte(frameValues, offset + 3)
                local w = string.byte(frameValues, offset + 4)

                LEDs.updateBuffer(j, g, r, b, w)
            end
            LEDs.writeBuffer()
        end

        fd:close();
        fd = nil
    end
end

function printDump(o)
    print(dump(o))
end

function dump(o)
    print("dump")
    print(o)
    if type(o) == 'table' then
        local s = '{ '
        for k, v in pairs(o) do
            if type(k) ~= 'number' then
                k = '"' .. k .. '"'
            end
            s = s .. '[' .. k .. '] = ' .. dump(v) .. ','
        end
        return s .. '} '
    else
        return tostring(o)
    end
end

function openUDPSocket()
    print("Start UDP")
    local udpSocket = net.createUDPSocket()
    udpSocket:dns("nodemcu-controller", function(conn, ip)
        print("Found controller ip " .. ip)
        udpSocket:on("receive", function(s, data, port, ip)
            --print(string.byte(data, 1, string.len(data)))

            if string.len(data) % 5 ~= 0 then
                udpSocket:send(1234, ip, "error: invalid data size")
            else
                for i = 0, string.len(data) / 5 - 1 do
                    local ledid = string.byte(data, i * 5 + 1)
                    local r = string.byte(data, i * 5 + 2)
                    local g = string.byte(data, i * 5 + 3)
                    local b = string.byte(data, i * 5 + 4)
                    local w = string.byte(data, i * 5 + 5)
                    --print(string.format("Set %d %d %d %d %d ", ledid, r, g, b, w))
                    LEDs.updateRGBW(ledid, g, r, b, w)
                end
            end
        end)
        udpSocket:send(1234, ip, "{id: 3, ledCount: 188, bytesPerLed: 4}")
    end)
end

if file.exists("OTA.update") then
    print("OTA file exists")
    file.remove("OTA.update")
    LFS.HTTP_OTA("sirmixalot", "/", "lfs.img")
else
    LEDs.init()
    setupWebServer()
end

