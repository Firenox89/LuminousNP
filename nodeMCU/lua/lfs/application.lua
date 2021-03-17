LEDs = require 'LED_control'

local nodeID, bytesPerLed, ledCount, segments, segmentList
local hearthbeatInterval = 10000

function setupWebServer()
    print("setup webserver")
    dofile('httpServer.lua')
    httpServer:listen(80)

    httpServer:use('/OTA', function(req, res)
        file.open("OTA.update", "w")
        file.close()
        node.task.post(function()
            node.restart()
        end)
        res:send(200)
    end)

    httpServer:use('/bar', function(req, res)
        res:send(200)
    end)

    httpServer:use('/baro', function(req, res)
        res:send(dump(readBarometer()))
    end)

    httpServer:use('/restart', function(req, res)
        LEDs.off(0)
        node.task.post(function()
            node.restart()
        end)
        res:send(200)
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
        local level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end
        LEDs.on(level)

        res:send(200)
    end)

    httpServer:use('/off', function(req, res)
        local level = tonumber(req.query.level)
        if level == nil then
            level = 0
        end

        LEDs.off(level)
        res:send(200)
    end)

    httpServer:use('/fill', function(req, res)
        local level = tonumber(req.query.level)
        local rawColor = req.query.color
        local r = tonumber(rawColor:sub(1, 2), 16)
        local g = tonumber(rawColor:sub(3, 4), 16)
        local b = tonumber(rawColor:sub(5, 6), 16)
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
        stopEffect()
        playEffect()
        res:send(200)
    end)

    httpServer:use('/startStream', function(req, res)
        openUDPSocket()
        res:send(200)
    end)

    httpServer:use('/startEffect', function(req, res)
        print("Start effect")

        stopEffect()
        LFS.HTTP_Download(
                "nodemcu-controller",
                "/",
                "effectFile",
                "?id=" .. nodeID,
                "current.effect",
                function()
                    print("effect file download complete")
                    playEffect()
                    res:send(200)
                end)
    end)

    httpServer:use('/updateConfig', function(req, res)
        local id = req.query.id
        local bytesPerLed = req.query.bytesperled
        local ledcount = req.query.ledcount
        local segments = req.query.segments

        saveConfig(id, bytesPerLed, ledcount, segments)
        loadConfig()
        res:send(200)
    end)

    registerAtController()
end

function loadConfig()
    if file.exists("node.config") then
        local fd = file.open("node.config", "r")
        nodeID = fd:readline():gsub("%s+", "")
        bytesPerLed = fd:readline():gsub("%s+", "")
        ledCount = fd:readline():gsub("%s+", "")
        segments = fd:readline():gsub("%s+", "")
        segmentList = splitString(segments, ",")
        fd:close()
        print("Config loaded")
        print("ID " .. nodeID)
        print("Bytes per LED " .. bytesPerLed)
        print("LED Count " .. ledCount)
        print("segments " .. segments)
    else
        print("No config files found")
        nodeID = "No ID set"
        bytesPerLed = 4
        ledCount = 0
        segments = ""
        segmentList = ""
    end
end

function saveConfig(id, bytesPerLed, ledcount, segments)
    local fd = file.open("node.config", "w")
    fd:writeline(id)
    fd:writeline(bytesPerLed)
    fd:writeline(ledcount)
    fd:writeline(segments)
    fd:close()
end

function splitString (inputstr, sep)
    if sep == nil then
        sep = "%s"
    end
    local t = {}
    for str in string.gmatch(inputstr, "([^" .. sep .. "]+)") do
        table.insert(t, str)
    end
    return t
end

function registerAtController()
    tmr.create():alarm(hearthbeatInterval, tmr.ALARM_AUTO, function()
        sendConfigToController()
    end)
end

function sendConfigToController()
    http.post('http://nodemcu-controller/register',
            'Content-Type: application/json\r\n',
            '{"id": "' .. nodeID ..
                    '", "ledCount": ' .. ledCount ..
                    ', "bytesPerLed": ' .. bytesPerLed ..
                    ', "segments":[' .. segments .. ']}',
            function(code, data)
                if (code < 0) then
                    print("Register at controller failed. code " .. code)
                end
            end)
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

busId = 0
sdaGPIO = 4
sclGPIO = 0

i2c.setup(busId, sdaGPIO, sclGPIO, i2c.FAST)

function readBarometer()
    deviceId = 59
    startRegister = 247 --0xF7
    readLength = 7

    i2c.start(busId)
    i2c.address(busId, deviceId, i2c.TRANSMITTER)
    i2c.write(busId, startRegister)
    i2c.stop(busId)

    i2c.start(busId)
    i2c.address(busId, deviceId, i2c.RECEIVER)
    local data = i2c.read(busId, readLength)
    i2c.stop(busId)

    local press_msb = string.byte(data, 0)
    local press_lsb = string.byte(data, 1)
    local press_xlsb = string.byte(data, 2)

    local temp_msb = string.byte(data, 0)
    local temp_lsb = string.byte(data, 1)
    local temp_xlsb = string.byte(data, 2)
    printDump(c)
end

if file.exists("OTA.update") then
    print("OTA file exists")
    file.remove("OTA.update")
    LFS.HTTP_Download("nodemcu-controller", "/", "lfs.img", "", "lfs.img", function()
        wifi.setmode(wifi.NULLMODE, false)
        collectgarbage();
        collectgarbage()
        -- run as separate task to maximise RAM available
        node.task.post(function()
            node.LFS.reload("lfs.img")
        end)
    end)
else
    loadConfig()
    setupWebServer()
end

