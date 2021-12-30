function handle (eventData)
    if eventData.geboorteDatum == nil then
        error("test: geboorteDatum is not present.")
    end

    local age = os.time() - eventData.geboorteDatum

    local t = os.date("*t", os.time())
    t.year = t.year - 18
    local ago = os.time() - os.time(t)

    local isOver18 = age > ago

    if not isOver18 then
        error("test: age is not over 18")
    end
end
