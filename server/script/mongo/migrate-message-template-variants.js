// Run with:
// mongosh <connection-string> --file migrate-message-template-variants.js

const collection = db.getCollection("message_templates");

if (collection.getIndexes().some((index) => index.name === "unique_name")) {
    collection.dropIndex("unique_name");
}

collection.find({ type: { $exists: false } }).forEach((template) => {
    collection.updateOne(
        { _id: template._id },
        {
            $set: {
                type: template.name,
                name: "默认模板",
                is_default: true,
                updated_at: new Date()
            }
        }
    );
});

const variablesByType = {
    "user-comment-approval": ["PostURL"],
    "user-comment-disapproval": ["PostURL", "Reason"],
    "user-comment-reply": ["PostURL"],
    "friend-approval": ["FriendPageURL"],
    "friend-rejection": ["FriendPageURL", "Reason"]
};

collection.find({ content: /%s/ }).forEach((template) => {
    let content = template.content;
    (variablesByType[template.type] || []).forEach((variable) => {
        content = content.replace("%s", `{{.${variable}}}`);
    });
    if (content.includes("%s")) {
        throw new Error(`Unsupported legacy placeholders in message template ${template._id}`);
    }
    collection.updateOne(
        { _id: template._id },
        { $set: { content: content, updated_at: new Date() } }
    );
});

collection.createIndex(
    { type: 1, name: 1 },
    { name: "unique_type_name", unique: true }
);
collection.createIndex(
    { type: 1, is_default: 1 },
    {
        name: "unique_default_per_type",
        unique: true,
        partialFilterExpression: { is_default: true }
    }
);
