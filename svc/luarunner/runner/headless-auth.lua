function handle(eventData)
    local jsonModelScratch = [[
{
  "r": {
    "m": {
      "getWlzIndicatieVoorIndicatieID": "#W"
    }
  },
  "W": {
    "m": { "grondslag": "#G", "geindiceerdZorgzwaartepakket": "#GZ", "geindiceerdeFunctie": "#GF", "beperking": "#B", "stoornis": "#S", "stoornisScore": "#SS", "wzd": "#WZD" },
    "f": ["wlzindicatieID", "bsn", "besluitnummer", "soortWlzindicatie", "afgiftedatum", "ingangsdatum", "einddatum", "meerzorg", "commentaar"],
    "p": {
      "wlzindicatieID": { $$ReplaceWithRecordID$$ }
    }
  },
  "G": {
    "f" : ["id", "grondslag", "volgorde", "wlzindicatieID"]
  },
  "GZ": {
    "f" : ["id", "zzpCode", "ingangsdatum", "einddatum", "klasse", "voorkeurClient", "instellingVoorkeur", "financiering", "commentaar", "wlzindicatieID"]
  },
  "GF": {
    "f": ["id", "functiecode", "ingangsdatum", "einddatum", "klasse", "opslag", "leveringsVoorwaarde", "vervoer", "instellingVoorkeur", "financiering", "commentaar", "wlzindicatieID"]
  },
  "B": {
    "m": {"beperkingScores":  "#BS"},
    "f": ["id", "categorie", "duur", "commentaar", "wlzindicatieID"]
  },
  "S": {
    "f": ["id", "beperkingVraag", "beperkingScore", "commentaar", "beperkingID"]
  },
  "SS": {
    "f": ["id", "stoornisVraag", "stoornisScore", "commentaar", "wlzindicatieID"]
  },
  "WZD": {
    "f": ["wzdVerklaring", "ingangsdatum", "einddatum"]
  },

  "BS": {
    "f": ["id", "beperkingVraag", "beperkingScore", "commentaar", "beperkingID"]
  }
}]]

    local jsonModel = jsonModelScratch:gsub("%$%$ReplaceWithRecordID%$%$", eventData.wlzindicatieID)
    local modelPath = "/tst/netwerkmodel/v1/graphql"
    if not headlessAuthorization(eventData.clientID, eventData.redirectID, eventData.audience, jsonModel, modelPath) then
        error("lua headless authorization")
    end
end