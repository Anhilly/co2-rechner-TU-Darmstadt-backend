console.log("Initiating MongoDB");

//db.createCollection('sample_collection');

db.createUser(
    {
        user: "backend",
        pwd: "test1234",
        roles: [
            {
                role: "readWrite",
                db: "co2Rechner"
            }
        ]
    }
);