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
            LEDs.fill(level, g, r, b, w)
        end
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
        r = tonumber(rawColor:sub(1, 2),16)
        g = tonumber(rawColor:sub(3, 4),16)
        b = tonumber(rawColor:sub(5, 6),16)
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
end

function printDump(o)
    print(dump(o))
end

function dump(o)
    print("dump")
    print(o)
    if type(o) == 'table' then
        local s = '{ '
        for k,v in pairs(o) do
            if type(k) ~= 'number' then k = '"'..k..'"' end
            s = s .. '['..k..'] = ' .. dump(v) .. ','
        end
        return s .. '} '
    else
        return tostring(o)
    end
end

if file.exists("OTA.update") then
    print("OTA file exists")
    file.remove("OTA.update")
    LFS.HTTP_OTA("sirmixalot", "/", "lfs.img")
else
    LEDs.init()
    setupWebServer()
end

