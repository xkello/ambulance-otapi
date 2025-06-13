const mongoHost = process.env.HOSPITAL_API_MONGODB_HOST
const mongoPort = process.env.HOSPITAL_API_MONGODB_PORT

const mongoUser = process.env.HOSPITAL_API_MONGODB_USERNAME
const mongoPassword = process.env.HOSPITAL_API_MONGODB_PASSWORD

const database = process.env.HOSPITAL_API_MONGODB_DATABASE
const collection = process.env.HOSPITAL_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
        print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })

//insert sample data
let result = db[collection].insertMany([
    {
        "id": "hospital-ba",
        "name": "Hospital Bratislava",
        "address": "123",
        "predefinedRoles": [
            { "value": "Doctor", "code": "rhinitis" },
            { "value": "Nurse", "code": "checkup" },
            { "value": "Transporter", "code": "jason-statham" }
        ]
    },
    {
        "id": "hospital-nr",
        "name": "Hospital Nitra",
        "address": "321",
        "predefinedRoles": [
            { "value": "Doctor", "code": "rhinitis" },
            { "value": "Nurse", "code": "checkup" },
            { "value": "Transporter", "code": "jason-statham" }
        ]
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);
