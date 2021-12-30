local data = graphql("https://gqlciz.staging.n-id.network/gql", [[
    query ClientData($id: UUID!) {
        client(id: $id) {
            geheimeClient
            director {
                name
                endpoint
            }
            person {
                geboorteDatum
                geslacht
                geslachtsnaam
                voorletters
                voorvoegselGeslachtnaam
            }
        }
    }
]], {id = "6956bf5b-eb62-4d74-8de4-3587f4c3f42b"})

local age = os.time() - data.client.person.geboorteDatum
local isOver18 = age > os.time{year=18}

return isOver18
