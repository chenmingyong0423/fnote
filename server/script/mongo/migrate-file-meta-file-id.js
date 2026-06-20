const collection = db.getCollection("file_meta");
let migrated = 0;

collection.find({ file_id: { $type: "object" } }).forEach((document) => {
    if (typeof document.file_id?.Data !== "string") {
        throw new Error(`invalid legacy file_id for document ${document._id}`);
    }
    collection.updateOne(
        { _id: document._id },
        { $set: { file_id: Binary.createFromBase64(document.file_id.Data, 0) } },
    );
    migrated++;
});

print(`migrated ${migrated} file_meta.file_id values`);
