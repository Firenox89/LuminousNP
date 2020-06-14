local host, dir, fileuri, query, localfile, callback = ...

local doRequest, firstRec, subsRec, finalise
local n, total, size = 0, 0

doRequest = function(socket, hostIP) -- luacheck: no unused
    if hostIP then
        local con = net.createConnection(net.TCP,0)
        -- Note that the current dev version can only accept uncompressed LFS images
        con:on("connection",function(sck)
            print("on connect")
            local request = table.concat( {
                    "GET "..dir..fileuri..query.." HTTP/1.1",
                    "User-Agent: ESP8266 app (linux-gnu)",
                    "Accept: application/octet-stream",
                    "Accept-Encoding: identity",
                    "Host: "..host,
                    "Connection: close",
                "", "", }, "\r\n")
            print(request)
            sck:send(request)
            sck:on("receive",firstRec)
        end)
        print("connect")
        con:connect(80,hostIP)
    end
end

firstRec = function (sck,rec)
    -- Process the headers; only interested in content length
    local i      = rec:find('\r\n\r\n',1,true) or 1
    local header = rec:sub(1,i+1):lower()
    size         = tonumber(header:match('\ncontent%-length: *(%d+)\r') or 0)
    print(rec:sub(1, i+1))
    if size > 0 then
        sck:on("receive",subsRec)
        file.open(localfile, 'w')
        subsRec(sck, rec:sub(i+4))
    else
        sck:on("receive", nil)
        sck:close()
        print("GET failed")
    end
end

subsRec = function(sck,rec)
    total, n = total + #rec, n + 1
    if n % 4 == 1 then
        sck:hold()
        node.task.post(0, function() sck:unhold() end)
    end
    uart.write(0,('%u of %u, '):format(total, size))
    file.write(rec)
    if total == size then finalise(sck) end
end

finalise = function(sck)
    file.close()
    sck:on("receive", nil)
    sck:close()
    print("Download finished")
    local s = file.stat(localfile)
    if (s and size == s.size) then
        callback()
    else
        if (s) then
            print"File not saved"
        else
            print("Size mismatch got " .. s.size .. ", needed " .. size)
        end
    end
end

net.dns.resolve(host, doRequest)
